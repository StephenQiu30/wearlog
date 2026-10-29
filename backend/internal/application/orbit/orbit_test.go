package orbit

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"

	accountapp "github.com/StephenQiu30/wearlog-server/backend/internal/application/account"
)

type mockAuthenticator struct {
	user accountapp.User
	err  error
}

func (m *mockAuthenticator) CurrentUser(ctx context.Context, session string) (accountapp.User, error) {
	if m.err != nil {
		return accountapp.User{}, m.err
	}
	return m.user, nil
}

type mockAnalyzer struct {
	result *AnalysisResult
	err    error
}

func (m *mockAnalyzer) Analyze(ctx context.Context, photoURL string) (*AnalysisResult, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.result, nil
}

type mockMatting struct {
	items []GarmentItem
	err   error
}

func (m *mockMatting) MatGarments(ctx context.Context, photoURL string, candidates []GarmentCandidate) ([]GarmentItem, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.items, nil
}

type mockTryOn struct {
	resultURL string
	err       error
}

func (m *mockTryOn) TryOn(ctx context.Context, modelPhotoURL string, garments []GarmentItem) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.resultURL, nil
}

type mockOrbitGen struct {
	sheet *SpriteSheetResource
	err   error
}

func (m *mockOrbitGen) Generate360(ctx context.Context, photoURL string) (*SpriteSheetResource, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.sheet, nil
}

type mockRepo struct {
	mu    sync.RWMutex
	looks map[string]*Look
}

func newMockRepo() *mockRepo {
	return &mockRepo{looks: make(map[string]*Look)}
}

func (m *mockRepo) SaveLook(ctx context.Context, look *Look) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.looks[look.ID] = look
	return nil
}

func (m *mockRepo) GetLook(ctx context.Context, lookID string) (*Look, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	l, ok := m.looks[lookID]
	if !ok {
		return nil, nil
	}
	return l, nil
}

func (m *mockRepo) ListLooksByUserID(ctx context.Context, userID string, limit, offset int) ([]*Look, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var res []*Look
	for _, l := range m.looks {
		if l.UserID == userID {
			res = append(res, l)
		}
	}
	return res, nil
}

func TestOrbitServiceAnalyzeOutfit(t *testing.T) {
	ctx := context.Background()
	analyzer := &mockAnalyzer{
		result: &AnalysisResult{
			Season:    SeasonAutumn,
			MinTemp:   15,
			MaxTemp:   22,
			StyleTags: []string{"casual", "minimalist"},
			GarmentPrompts: []GarmentCandidate{
				{
					Category:     CategoryTop,
					Name:         "White Cotton T-Shirt",
					Subcategory:  "t_shirt",
					BoundingBox:  BoundingBox{0.2, 0.3, 0.5, 0.7},
					ColorPalette: []string{"#FFFFFF"},
					Material:     "cotton",
					Confidence:   0.95,
				},
			},
		},
	}
	matting := &mockMatting{
		items: []GarmentItem{
			{
				ID:            "g-1",
				Category:      CategoryTop,
				Name:          "White Cotton T-Shirt",
				Subcategory:   "t_shirt",
				MattingPNGURL: "https://cdn.example.com/matting/g-1.png",
				BoundingBox:   BoundingBox{0.2, 0.3, 0.5, 0.7},
				ColorPalette:  []string{"#FFFFFF"},
				Material:      "cotton",
				Confidence:    0.95,
			},
		},
	}
	tryOn := &mockTryOn{resultURL: "https://cdn.example.com/tryon/res.webp"}
	orbitGen := &mockOrbitGen{
		sheet: &SpriteSheetResource{
			SheetURL:    "https://cdn.example.com/orbit/sheet.webp",
			FramesCount: 16,
			LayoutRows:  4,
			LayoutCols:  4,
			FrameWidth:  512,
			FrameHeight: 683,
			TotalSizeKB: 1120,
		},
	}
	repo := newMockRepo()
	auth := &mockAuthenticator{user: accountapp.User{ID: "user-123"}}

	svc, err := NewService(auth, analyzer, matting, tryOn, orbitGen, repo)
	if err != nil {
		t.Fatalf("unexpected error initializing service: %v", err)
	}

	// 1. Invalid Photo URL
	_, err = svc.AnalyzeOutfit(ctx, "session-123", "not-a-url")
	if !errors.Is(err, ErrEmptyPhotoURL) {
		t.Fatalf("expected ErrEmptyPhotoURL, got: %v", err)
	}

	// 2. Successful analysis and persistence
	look, err := svc.AnalyzeOutfit(ctx, "session-123", "https://example.com/outfit.jpg")
	if err != nil {
		t.Fatalf("AnalyzeOutfit failed: %v", err)
	}
	if look.Season != SeasonAutumn {
		t.Fatalf("expected SeasonAutumn, got: %s", look.Season)
	}
	if len(look.Garments) != 1 || look.Garments[0].ID != "g-1" {
		t.Fatalf("expected 1 garment with ID g-1, got: %v", look.Garments)
	}

	// 3. Try-On
	tryOnRes, err := svc.TryOn(ctx, "session-123", look.ID, look.Garments)
	if err != nil {
		t.Fatalf("TryOn failed: %v", err)
	}
	if tryOnRes.ResultImageURL != "https://cdn.example.com/tryon/res.webp" {
		t.Fatalf("unexpected try-on URL: %s", tryOnRes.ResultImageURL)
	}

	// 4. Generate 360°
	orbitRes, err := svc.Generate360(ctx, "session-123", look.ID)
	if err != nil {
		t.Fatalf("Generate360 failed: %v", err)
	}
	if orbitRes.SpriteSheet.FramesCount != 16 {
		t.Fatalf("expected 16 frames, got: %d", orbitRes.SpriteSheet.FramesCount)
	}

	// 5. Verify look in repo was updated with sprite sheet
	updatedLook, err := svc.GetLook(ctx, "session-123", look.ID)
	if err != nil {
		t.Fatalf("GetLook failed: %v", err)
	}
	if updatedLook.SpriteSheet == nil || updatedLook.SpriteSheet.SheetURL != orbitRes.SpriteSheet.SheetURL {
		t.Fatalf("Look was not updated with sprite sheet in repo")
	}
}

func TestOrbitServiceDependencyValidation(t *testing.T) {
	_, err := NewService(nil, nil, nil, nil, nil, nil)
	if err == nil {
		t.Fatal("expected error when passing nil dependencies")
	}
}

func TestOrbitServiceGoogleGoErrorFormatting(t *testing.T) {
	auth := &mockAuthenticator{user: accountapp.User{ID: "user-123"}}
	svc, _ := NewService(auth, &mockAnalyzer{}, &mockMatting{}, &mockTryOn{}, &mockOrbitGen{}, newMockRepo())
	_, err := svc.TryOn(context.Background(), "session-123", "look-1", nil)
	if !errors.Is(err, ErrEmptyGarments) {
		t.Fatalf("expected ErrEmptyGarments, got %v", err)
	}
	if strings.Contains(err.Error(), ".") {
		t.Fatalf("error message %q violates Google Go Style Guide (should not end with period)", err.Error())
	}
}
