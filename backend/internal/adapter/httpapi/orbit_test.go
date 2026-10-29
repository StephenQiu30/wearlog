package httpapi

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/StephenQiu30/wearlog-server/backend/internal/application/orbit"
)

type orbitTransportStub struct {
	look *orbit.Look
	err  error
}

func (s *orbitTransportStub) AnalyzeOutfit(ctx context.Context, session, photoURL string) (*orbit.Look, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.look, nil
}

func (s *orbitTransportStub) TryOn(ctx context.Context, session, lookID string, garmentsToWear []orbit.GarmentItem) (*orbit.TryOnResult, error) {
	if s.err != nil {
		return nil, s.err
	}
	return &orbit.TryOnResult{
		LookID:             lookID,
		ResultImageURL:     "https://cdn.example.com/tryon-result.webp",
		SelectedGarmentIDs: []string{"g-1"},
		LatencyMS:          1200,
	}, nil
}

func (s *orbitTransportStub) Generate360(ctx context.Context, session, lookID string) (*orbit.Orbit360Result, error) {
	if s.err != nil {
		return nil, s.err
	}
	return &orbit.Orbit360Result{
		LookID: lookID,
		SpriteSheet: orbit.SpriteSheetResource{
			SheetURL:    "https://cdn.example.com/orbit-spritesheet.webp",
			FramesCount: 16,
			LayoutRows:  4,
			LayoutCols:  4,
			FrameWidth:  512,
			FrameHeight: 683,
			TotalSizeKB: 1120,
		},
		LatencyMS: 8500,
	}, nil
}

func (s *orbitTransportStub) GetLook(ctx context.Context, session, lookID string) (*orbit.Look, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.look, nil
}

func orbitFixture() *orbit.Look {
	now := time.Date(2026, 10, 29, 12, 0, 0, 0, time.UTC)
	return &orbit.Look{
		ID:               "11111111-1111-4111-8111-111111111111",
		UserID:           "22222222-2222-4222-8222-222222222222",
		OriginalPhotoURL: "https://example.com/photo.jpg",
		DepthMapURL:      "https://example.com/depth.png",
		Season:           orbit.SeasonAutumn,
		MinTemp:          15,
		MaxTemp:          22,
		StyleTags:        []string{"casual", "city_boy"},
		Garments: []orbit.GarmentItem{
			{
				ID:            "g-1",
				Category:      orbit.CategoryTop,
				Name:          "Cream Knit Sweater",
				Subcategory:   "knitwear",
				MattingPNGURL: "https://example.com/g-1.png",
				BoundingBox:   orbit.BoundingBox{0.2, 0.25, 0.55, 0.75},
				ColorPalette:  []string{"#FDFBF7"},
				Material:      "wool",
				Confidence:    0.95,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func orbitRouter(t *testing.T, service OrbitHTTPService) *Router {
	t.Helper()
	probe := probeFunc(func(context.Context) error { return nil })
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	router, err := NewRouterWithOrbit(
		context.Background(), false, probe,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		NewOrbitHandler(service, true), time.Second, logger,
	)
	if err != nil {
		t.Fatal(err)
	}
	return router
}

func TestOrbitAnalyzeEndpoint(t *testing.T) {
	router := orbitRouter(t, &orbitTransportStub{look: orbitFixture()})

	// 1. Unauthenticated request returns 401
	unauthReq := httptest.NewRequest(http.MethodPost, "/orbit/analyze", strings.NewReader(`{"photo_url":"https://example.com/photo.jpg"}`))
	unauthReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, unauthReq)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}

	// 2. Authenticated request returns 200 with LookResponse
	authReq := httptest.NewRequest(http.MethodPost, "/orbit/analyze", strings.NewReader(`{"photo_url":"https://example.com/photo.jpg"}`))
	authReq.Header.Set("Content-Type", "application/json")
	authReq.AddCookie(&http.Cookie{Name: sessionCookieName, Value: "synthetic-session-token"})
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, authReq)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "Cream Knit Sweater") {
		t.Fatalf("response missing garment name: %s", rec.Body.String())
	}
}

func TestOrbitTryOnEndpoint(t *testing.T) {
	router := orbitRouter(t, &orbitTransportStub{look: orbitFixture()})

	body := `{"selected_garment_ids":["g-1"],"garments":[{"id":"g-1","category":"top","name":"Cream Knit","subcategory":"knitwear","matting_png_url":"https://example.com/g-1.png","bounding_box":[0.2,0.25,0.55,0.75],"color_palette":["#FFFFFF"],"material":"wool","confidence":0.95}]}`
	req := httptest.NewRequest(http.MethodPost, "/orbit/looks/11111111-1111-4111-8111-111111111111/try-on", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: sessionCookieName, Value: "synthetic-session-token"})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "tryon-result.webp") {
		t.Fatalf("response missing try-on result: %s", rec.Body.String())
	}
}

func TestOrbitGenerate360Endpoint(t *testing.T) {
	router := orbitRouter(t, &orbitTransportStub{look: orbitFixture()})

	req := httptest.NewRequest(http.MethodPost, "/orbit/looks/11111111-1111-4111-8111-111111111111/generate-360", nil)
	req.AddCookie(&http.Cookie{Name: sessionCookieName, Value: "synthetic-session-token"})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "orbit-spritesheet.webp") {
		t.Fatalf("response missing sprite sheet URL: %s", rec.Body.String())
	}
}
