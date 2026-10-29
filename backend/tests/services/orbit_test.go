//go:build services

package services

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"testing"
	"time"

	store "github.com/StephenQiu30/wearlog-server/backend/internal/adapter/postgres"
	accountapp "github.com/StephenQiu30/wearlog-server/backend/internal/application/account"
	orbitapp "github.com/StephenQiu30/wearlog-server/backend/internal/application/orbit"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type mockFashionAnalyzer struct{}

func (m *mockFashionAnalyzer) Analyze(ctx context.Context, photoURL string) (*orbitapp.AnalysisResult, error) {
	return &orbitapp.AnalysisResult{
		Season:    orbitapp.SeasonAutumn,
		MinTemp:   15,
		MaxTemp:   22,
		StyleTags: []string{"casual", "city_boy"},
		Candidates: []orbitapp.GarmentCandidate{
			{
				Category:    orbitapp.CategoryTop,
				Name:        "Mock Sweater",
				Subcategory: "knitwear",
				BoundingBox: orbitapp.BoundingBox{0.2, 0.25, 0.55, 0.75},
				Confidence:  0.95,
			},
		},
	}, nil
}

type mockMattingProvider struct{}

func (m *mockMattingProvider) MatGarments(ctx context.Context, photoURL string, candidates []orbitapp.GarmentCandidate) ([]orbitapp.GarmentItem, error) {
	items := make([]orbitapp.GarmentItem, 0, len(candidates))
	for _, c := range candidates {
		items = append(items, orbitapp.GarmentItem{
			ID:            uuid.NewString(),
			Category:      c.Category,
			Name:          c.Name,
			Subcategory:   c.Subcategory,
			MattingPNGURL: "https://example.com/matting.png",
			BoundingBox:   c.BoundingBox,
			ColorPalette:  []string{"#FFFFFF"},
			Material:      "cotton",
			Confidence:    c.Confidence,
		})
	}
	return items, nil
}

type mockTryOnProvider struct{}

func (m *mockTryOnProvider) TryOn(ctx context.Context, modelPhotoURL string, garments []orbitapp.GarmentItem) (string, error) {
	return "https://example.com/tryon-result.webp", nil
}

type mockOrbitGenerator struct{}

func (m *mockOrbitGenerator) Generate360(ctx context.Context, photoURL string) (*orbitapp.SpriteSheetResource, error) {
	return &orbitapp.SpriteSheetResource{
		SheetURL:    "https://example.com/spritesheet.webp",
		FramesCount: 16,
		LayoutRows:  4,
		LayoutCols:  4,
		FrameWidth:  512,
		FrameHeight: 683,
		TotalSizeKB: 1024,
	}, nil
}

func TestOrbitPersistenceLifecycle(t *testing.T) {
	environment := loadServiceEnvironment(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	base, err := gorm.Open(postgres.Open(environment.databaseURL), &gorm.Config{Logger: logger.Discard})
	serviceOK(t, "open orbit test database", err)
	schema := strings.ReplaceAll(serviceID(t), "-", "_")
	serviceOK(t, "create orbit test schema", base.WithContext(ctx).Exec("CREATE SCHEMA "+schema).Error)
	t.Cleanup(func() {
		cleanup, done := context.WithTimeout(context.Background(), 5*time.Second)
		defer done()
		if err := base.WithContext(cleanup).Exec("DROP SCHEMA " + schema + " CASCADE").Error; err != nil {
			t.Error("orbit schema cleanup failed")
		}
	})

	databaseURL, err := url.Parse(environment.databaseURL)
	serviceOK(t, "parse orbit test database", err)
	query := databaseURL.Query()
	query.Set("search_path", schema)
	databaseURL.RawQuery = query.Encode()
	database, err := gorm.Open(postgres.Open(databaseURL.String()), &gorm.Config{Logger: logger.Discard})
	serviceOK(t, "open isolated orbit schema", err)
	serviceOK(t, "migrate orbit schema", store.Migrate(ctx, database))

	accounts, err := accountapp.NewAccountService(store.NewAccountRepository(database))
	serviceOK(t, "construct orbit account service", err)
	firstUser, err := accounts.Register(ctx, accountapp.RegisterAccountInput{
		Email:       "orbit-first@example.test",
		DisplayName: "Orbit User 1",
		Password:    "correct-password-first",
	})
	serviceOK(t, "register first orbit owner", err)

	secondUser, err := accounts.Register(ctx, accountapp.RegisterAccountInput{
		Email:       "orbit-second@example.test",
		DisplayName: "Orbit User 2",
		Password:    "correct-password-second",
	})
	serviceOK(t, "register second orbit owner", err)

	repo := store.NewOrbitRepository(database)
	lookID := uuid.NewString()
	now := time.Now().UTC()
	look := &orbitapp.Look{
		ID:               lookID,
		UserID:           firstUser.User.ID,
		OriginalPhotoURL: "https://cdn.example.com/outfit1.jpg",
		DepthMapURL:      "https://cdn.example.com/outfit1_depth.png",
		Season:           orbitapp.SeasonAutumn,
		MinTemp:          12,
		MaxTemp:          20,
		StyleTags:        []string{"cleanfit", "normcore"},
		Garments: []orbitapp.GarmentItem{
			{
				ID:            uuid.NewString(),
				Category:      orbitapp.CategoryTop,
				Name:          "Wool Sweater",
				Subcategory:   "knitwear",
				MattingPNGURL: "https://cdn.example.com/sweater.png",
				BoundingBox:   orbitapp.BoundingBox{0.2, 0.25, 0.55, 0.75},
				ColorPalette:  []string{"#E5E0D8"},
				Material:      "wool",
				Confidence:    0.96,
			},
			{
				ID:            uuid.NewString(),
				Category:      orbitapp.CategoryBottom,
				Name:          "Wide Trousers",
				Subcategory:   "slacks",
				MattingPNGURL: "https://cdn.example.com/pants.png",
				BoundingBox:   orbitapp.BoundingBox{0.5, 0.2, 0.85, 0.8},
				ColorPalette:  []string{"#2B2B2B"},
				Material:      "cotton",
				Confidence:    0.94,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// 1. Save Look
	serviceOK(t, "save orbit look", repo.SaveLook(ctx, look))

	// 2. Get Look
	retrieved, err := repo.GetLook(ctx, lookID)
	serviceOK(t, "get orbit look", err)
	if retrieved.ID != lookID || retrieved.UserID != firstUser.User.ID {
		t.Fatalf("unexpected look: %+v", retrieved)
	}
	if len(retrieved.Garments) != 2 {
		t.Fatalf("expected 2 garments, got %d", len(retrieved.Garments))
	}
	if retrieved.Season != orbitapp.SeasonAutumn {
		t.Fatalf("expected Autumn season, got %s", retrieved.Season)
	}

	// 3. Update Look with SpriteSheet
	retrieved.SpriteSheet = &orbitapp.SpriteSheetResource{
		SheetURL:    "https://cdn.example.com/spritesheet.webp",
		FramesCount: 16,
		LayoutRows:  4,
		LayoutCols:  4,
		FrameWidth:  512,
		FrameHeight: 683,
		TotalSizeKB: 1024,
	}
	serviceOK(t, "update orbit look with spritesheet", repo.SaveLook(ctx, retrieved))
	updated, err := repo.GetLook(ctx, lookID)
	serviceOK(t, "get updated orbit look", err)
	if updated.SpriteSheet == nil || updated.SpriteSheet.FramesCount != 16 {
		t.Fatalf("expected sprite sheet to be persisted, got: %+v", updated.SpriteSheet)
	}

	// 4. List Looks by User ID
	userLooks, err := repo.ListLooksByUserID(ctx, firstUser.User.ID, 10, 0)
	serviceOK(t, "list looks by user id", err)
	if len(userLooks) != 1 {
		t.Fatalf("expected 1 look for first user, got %d", len(userLooks))
	}
	secondUserLooks, err := repo.ListLooksByUserID(ctx, secondUser.User.ID, 10, 0)
	serviceOK(t, "list looks for second user", err)
	if len(secondUserLooks) != 0 {
		t.Fatalf("expected 0 looks for second user, got %d", len(secondUserLooks))
	}

	// 5. Look not found
	_, err = repo.GetLook(ctx, uuid.NewString())
	if !errors.Is(err, orbitapp.ErrLookNotFound) {
		t.Fatalf("expected ErrLookNotFound, got %v", err)
	}

	// 6. Service level isolation and access control
	service, err := orbitapp.NewService(
		repo,
		accounts,
		&mockFashionAnalyzer{},
		&mockMattingProvider{},
		&mockTryOnProvider{},
		&mockOrbitGenerator{},
	)
	serviceOK(t, "create orbit app service", err)

	// Owner can read their own look
	ownerLook, err := service.GetLook(ctx, firstUser.Token, lookID)
	serviceOK(t, "get look via service for owner", err)
	if ownerLook.ID != lookID {
		t.Fatalf("unexpected look id %s", ownerLook.ID)
	}

	// Other user cannot read the first user's look (privacy boundary: returns ErrLookNotFound)
	_, err = service.GetLook(ctx, secondUser.Token, lookID)
	if !errors.Is(err, orbitapp.ErrLookNotFound) {
		t.Fatalf("expected ErrLookNotFound for unauthorized user, got %v", err)
	}
}
