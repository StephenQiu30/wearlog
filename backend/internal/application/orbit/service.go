package orbit

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

// Service encapsulates the business logic for analyzing outfits, segmenting
// individual garments into digital wardrobe items, and orchestrating 360° and try-on operations.
type Service struct {
	authenticator Authenticator
	analyzer      FashionAnalyzer
	matting       MattingProvider
	tryOn         TryOnProvider
	orbitGen      OrbitGenerator
	repo          Repository
	now           func() time.Time
}

// NewService validates dependencies and returns a configured orbit Service.
func NewService(
	authenticator Authenticator,
	analyzer FashionAnalyzer,
	matting MattingProvider,
	tryOn TryOnProvider,
	orbitGen OrbitGenerator,
	repo Repository,
) (*Service, error) {
	if authenticator == nil || analyzer == nil || matting == nil || tryOn == nil || orbitGen == nil || repo == nil {
		return nil, errors.New("orbit dependencies cannot be nil")
	}
	return &Service{
		authenticator: authenticator,
		analyzer:      analyzer,
		matting:       matting,
		tryOn:         tryOn,
		orbitGen:      orbitGen,
		repo:          repo,
		now:           time.Now,
	}, nil
}

// AnalyzeOutfit processes an uploaded full-body photo using multi-modal LLM reasoning
// and alpha-matting segmentation, returning a persisted Look aggregate root.
func (s *Service) AnalyzeOutfit(ctx context.Context, session, photoURL string) (*Look, error) {
	user, err := s.authenticator.CurrentUser(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("authenticate user: %w", err)
	}
	if !isValidURL(photoURL) {
		return nil, ErrEmptyPhotoURL
	}

	analysis, err := s.analyzer.Analyze(ctx, photoURL)
	if err != nil {
		return nil, fmt.Errorf("fashion analysis failed: %w", err)
	}
	if analysis == nil {
		return nil, errors.New("fashion analysis returned empty result")
	}

	garments, err := s.matting.MatGarments(ctx, photoURL, analysis.GarmentPrompts)
	if err != nil {
		return nil, fmt.Errorf("garment matting failed: %w", err)
	}

	now := s.now().UTC()
	look := &Look{
		ID:               uuid.NewString(),
		UserID:           user.ID,
		OriginalPhotoURL: photoURL,
		Season:           analysis.Season,
		MinTemp:          analysis.MinTemp,
		MaxTemp:          analysis.MaxTemp,
		StyleTags:        analysis.StyleTags,
		Garments:         garments,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := s.repo.SaveLook(ctx, look); err != nil {
		return nil, fmt.Errorf("saving look failed: %w", err)
	}

	return look, nil
}

// TryOn replaces garments on the person in the look while preserving native pose.
func (s *Service) TryOn(ctx context.Context, session, lookID string, garmentsToWear []GarmentItem) (*TryOnResult, error) {
	user, err := s.authenticator.CurrentUser(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("authenticate user: %w", err)
	}
	if strings.TrimSpace(lookID) == "" {
		return nil, errors.New("look id cannot be empty")
	}
	if len(garmentsToWear) == 0 {
		return nil, ErrEmptyGarments
	}

	look, err := s.repo.GetLook(ctx, lookID)
	if err != nil {
		return nil, fmt.Errorf("fetch look failed: %w", err)
	}
	if look == nil || look.UserID != user.ID {
		return nil, ErrLookNotFound
	}

	started := s.now()
	resultImageURL, err := s.tryOn.TryOn(ctx, look.OriginalPhotoURL, garmentsToWear)
	if err != nil {
		return nil, fmt.Errorf("try-on provider failed: %w", err)
	}

	latency := s.now().Sub(started).Milliseconds()
	selectedIDs := make([]string, len(garmentsToWear))
	for i, g := range garmentsToWear {
		selectedIDs[i] = g.ID
	}

	return &TryOnResult{
		LookID:             lookID,
		ResultImageURL:     resultImageURL,
		SelectedGarmentIDs: selectedIDs,
		LatencyMS:          latency,
		UpdatedGarments:    garmentsToWear,
	}, nil
}

// Generate360 produces a 16-frame 4x4 WebP sprite sheet for 360° interactive rotation.
func (s *Service) Generate360(ctx context.Context, session, lookID string) (*Orbit360Result, error) {
	user, err := s.authenticator.CurrentUser(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("authenticate user: %w", err)
	}
	if strings.TrimSpace(lookID) == "" {
		return nil, errors.New("look id cannot be empty")
	}

	look, err := s.repo.GetLook(ctx, lookID)
	if err != nil {
		return nil, fmt.Errorf("fetch look failed: %w", err)
	}
	if look == nil || look.UserID != user.ID {
		return nil, ErrLookNotFound
	}

	started := s.now()
	sheet, err := s.orbitGen.Generate360(ctx, look.OriginalPhotoURL)
	if err != nil {
		return nil, fmt.Errorf("360 generation failed: %w", err)
	}

	look.SpriteSheet = sheet
	look.UpdatedAt = s.now().UTC()
	if err := s.repo.SaveLook(ctx, look); err != nil {
		return nil, fmt.Errorf("saving look sprite sheet failed: %w", err)
	}

	latency := s.now().Sub(started).Milliseconds()
	return &Orbit360Result{
		LookID:      lookID,
		SpriteSheet: *sheet,
		LatencyMS:   latency,
	}, nil
}

// GetLook retrieves a look by its identifier.
func (s *Service) GetLook(ctx context.Context, session, lookID string) (*Look, error) {
	user, err := s.authenticator.CurrentUser(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("authenticate user: %w", err)
	}
	if strings.TrimSpace(lookID) == "" {
		return nil, errors.New("look id cannot be empty")
	}
	look, err := s.repo.GetLook(ctx, lookID)
	if err != nil {
		return nil, fmt.Errorf("fetch look failed: %w", err)
	}
	if look == nil || look.UserID != user.ID {
		return nil, ErrLookNotFound
	}
	return look, nil
}

// BatchProcess concurrent demonstration showing errgroup usage conforming to Google Go Style Guide.
func (s *Service) BatchAnalyze(ctx context.Context, session string, photoURLs []string) ([]*Look, error) {
	g, gctx := errgroup.WithContext(ctx)
	results := make([]*Look, len(photoURLs))

	for i, u := range photoURLs {
		index, photoURL := i, u
		g.Go(func() error {
			look, err := s.AnalyzeOutfit(gctx, session, photoURL)
			if err != nil {
				return err
			}
			results[index] = look
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("batch analyze failed: %w", err)
	}
	return results, nil
}

func isValidURL(str string) bool {
	u, err := url.ParseRequestURI(str)
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}
