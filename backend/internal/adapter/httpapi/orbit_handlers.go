package httpapi

import (
	"context"
	"errors"
	"net/http"

	accountapp "github.com/StephenQiu30/wearlog-server/backend/internal/application/account"
	"github.com/StephenQiu30/wearlog-server/backend/internal/application/orbit"
	"github.com/danielgtaylor/huma/v2"
)

// OrbitHTTPService defines the contract required by HTTP handlers.
type OrbitHTTPService interface {
	AnalyzeOutfit(ctx context.Context, session, photoURL string) (*orbit.Look, error)
	TryOn(ctx context.Context, session, lookID string, garmentsToWear []orbit.GarmentItem) (*orbit.TryOnResult, error)
	Generate360(ctx context.Context, session, lookID string) (*orbit.Orbit360Result, error)
	GetLook(ctx context.Context, session, lookID string) (*orbit.Look, error)
}

// OrbitHandler bridges Huma HTTP operations and the Orbit application service.
type OrbitHandler struct {
	service      OrbitHTTPService
	secureCookie bool
}

// NewOrbitHandler creates an initialized OrbitHandler.
func NewOrbitHandler(service OrbitHTTPService, secureCookie bool) *OrbitHandler {
	return &OrbitHandler{service: service, secureCookie: secureCookie}
}

func (h *OrbitHandler) available(ctx context.Context, session string) error {
	if h == nil || h.service == nil {
		return huma.ErrorWithHeaders(newErrorResponse(http.StatusServiceUnavailable, requestID(ctx)), http.Header{"Retry-After": []string{"1"}})
	}
	if session == "" {
		return authenticatedSessionError(ctx, h.secureCookie)
	}
	return nil
}

func (h *OrbitHandler) analyze(ctx context.Context, input *analyzeOutfitInput) (*analyzeOutfitOutput, error) {
	if err := h.available(ctx, input.Session); err != nil {
		return nil, err
	}
	look, err := h.service.AnalyzeOutfit(ctx, input.Session, input.Body.PhotoURL)
	if err != nil {
		return nil, h.error(ctx, err)
	}
	return &analyzeOutfitOutput{
		RequestID: requestID(ctx),
		Body:      orbitLookResponse(look),
	}, nil
}

func (h *OrbitHandler) tryOn(ctx context.Context, input *tryOnInput) (*tryOnOutput, error) {
	if err := h.available(ctx, input.Session); err != nil {
		return nil, err
	}

	garments := make([]orbit.GarmentItem, len(input.Body.Garments))
	for i, g := range input.Body.Garments {
		garments[i] = orbit.GarmentItem{
			ID:            g.ID,
			Category:      orbit.GarmentCategory(g.Category),
			Name:          g.Name,
			Subcategory:   g.Subcategory,
			MattingPNGURL: g.MattingPNGURL,
			BoundingBox:   orbit.BoundingBox(g.BoundingBox),
			ColorPalette:  g.ColorPalette,
			Material:      g.Material,
			Confidence:    g.Confidence,
		}
	}

	result, err := h.service.TryOn(ctx, input.Session, input.LookID, garments)
	if err != nil {
		return nil, h.error(ctx, err)
	}

	return &tryOnOutput{
		RequestID: requestID(ctx),
		Body: TryOnResponse{
			LookID:             result.LookID,
			ResultImageURL:     result.ResultImageURL,
			SelectedGarmentIDs: result.SelectedGarmentIDs,
			LatencyMS:          result.LatencyMS,
		},
	}, nil
}

func (h *OrbitHandler) generate360(ctx context.Context, input *generate360Input) (*generate360Output, error) {
	if err := h.available(ctx, input.Session); err != nil {
		return nil, err
	}

	result, err := h.service.Generate360(ctx, input.Session, input.LookID)
	if err != nil {
		return nil, h.error(ctx, err)
	}

	return &generate360Output{
		RequestID: requestID(ctx),
		Body: Generate360Response{
			LookID:      result.LookID,
			SpriteSheet: orbitSpriteSheetDTO(result.SpriteSheet),
			LatencyMS:   result.LatencyMS,
		},
	}, nil
}

func (h *OrbitHandler) getLook(ctx context.Context, input *getLookInput) (*getLookOutput, error) {
	if err := h.available(ctx, input.Session); err != nil {
		return nil, err
	}

	look, err := h.service.GetLook(ctx, input.Session, input.LookID)
	if err != nil {
		return nil, h.error(ctx, err)
	}

	return &getLookOutput{
		RequestID: requestID(ctx),
		Body:      orbitLookResponse(look),
	}, nil
}

func (h *OrbitHandler) error(ctx context.Context, err error) error {
	switch {
	case errors.Is(err, accountapp.ErrAuthentication):
		return authenticatedSessionError(ctx, h.secureCookie)
	case errors.Is(err, orbit.ErrEmptyPhotoURL), errors.Is(err, orbit.ErrEmptyGarments), errors.Is(err, orbit.ErrInvalidCategory):
		return newErrorResponse(http.StatusBadRequest, requestID(ctx))
	case errors.Is(err, orbit.ErrLookNotFound):
		return newErrorResponse(http.StatusNotFound, requestID(ctx))
	case errors.Is(err, orbit.ErrProviderUnavailable), errors.Is(err, orbit.ErrTimeout):
		return huma.ErrorWithHeaders(newErrorResponse(http.StatusServiceUnavailable, requestID(ctx)), http.Header{"Retry-After": []string{"1"}})
	default:
		return newErrorResponse(http.StatusInternalServerError, requestID(ctx))
	}
}

func orbitLookResponse(look *orbit.Look) LookResponse {
	garmentDTOs := make([]GarmentItemDTO, len(look.Garments))
	for i, g := range look.Garments {
		garmentDTOs[i] = GarmentItemDTO{
			ID:            g.ID,
			Category:      string(g.Category),
			Name:          g.Name,
			Subcategory:   g.Subcategory,
			MattingPNGURL: g.MattingPNGURL,
			BoundingBox:   [4]float64(g.BoundingBox),
			ColorPalette:  g.ColorPalette,
			Material:      g.Material,
			Confidence:    g.Confidence,
		}
	}

	var sheetDTO *SpriteSheetDTO
	if look.SpriteSheet != nil {
		s := orbitSpriteSheetDTO(*look.SpriteSheet)
		sheetDTO = &s
	}

	return LookResponse{
		ID:               look.ID,
		OriginalPhotoURL: look.OriginalPhotoURL,
		DepthMapURL:      look.DepthMapURL,
		Season:           string(look.Season),
		MinTemp:          look.MinTemp,
		MaxTemp:          look.MaxTemp,
		StyleTags:        look.StyleTags,
		Garments:         garmentDTOs,
		SpriteSheet:      sheetDTO,
		CreatedAt:        look.CreatedAt,
		UpdatedAt:        look.UpdatedAt,
	}
}

func orbitSpriteSheetDTO(sheet orbit.SpriteSheetResource) SpriteSheetDTO {
	return SpriteSheetDTO{
		SheetURL:    sheet.SheetURL,
		FramesCount: sheet.FramesCount,
		LayoutRows:  sheet.LayoutRows,
		LayoutCols:  sheet.LayoutCols,
		FrameWidth:  sheet.FrameWidth,
		FrameHeight: sheet.FrameHeight,
		TotalSizeKB: sheet.TotalSizeKB,
	}
}
