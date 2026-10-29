// Package orbit implements multi-modal outfit deconstruction, client-side
// parallax coordination, and 360° keyframe sprite sheet generation according
// to the Woo product specification without self-hosted GPU compute.
package orbit

import (
	"errors"
	"time"
)

// Domain errors conforming to Google Go Style Guide (lowercase, no trailing punctuation).
var (
	ErrEmptyPhotoURL       = errors.New("photo url cannot be empty")
	ErrLookNotFound        = errors.New("look not found")
	ErrEmptyGarments       = errors.New("garments list cannot be empty")
	ErrInvalidCategory     = errors.New("invalid garment category")
	ErrProviderUnavailable = errors.New("external ai provider unavailable")
	ErrTimeout             = errors.New("orbit processing timed out")
	ErrInvalidBoundingBox  = errors.New("bounding box coordinates must be normalized between 0 and 1")
)

// Season describes the primary seasonal suitability classified by the multi-modal LLM.
type Season string

const (
	SeasonSpring    Season = "spring"
	SeasonSummer    Season = "summer"
	SeasonAutumn    Season = "autumn"
	SeasonWinter    Season = "winter"
	SeasonAllSeason Season = "all_season"
)

// GarmentCategory categorizes garments extracted from the outfit photo.
type GarmentCategory string

const (
	CategoryOuterwear GarmentCategory = "outerwear"
	CategoryTop       GarmentCategory = "top"
	CategoryBottom    GarmentCategory = "bottom"
	CategoryShoes     GarmentCategory = "shoes"
	CategoryAccessory GarmentCategory = "accessory"
)

// BoundingBox stores normalized [ymin, xmin, ymax, xmax] coordinates in [0, 1].
type BoundingBox [4]float64

// GarmentItem represents an individual segmented item in the digital wardrobe dock.
type GarmentItem struct {
	ID            string          `json:"id"`
	Category      GarmentCategory `json:"category"`
	Name          string          `json:"name"`
	Subcategory   string          `json:"subcategory"`
	MattingPNGURL string          `json:"matting_png_url"`
	BoundingBox   BoundingBox     `json:"bounding_box"`
	ColorPalette  []string        `json:"color_palette"`
	Material      string          `json:"material"`
	Confidence    float64         `json:"confidence"`
}

// SpriteSheetResource points to the 16-frame 4x4 WebP atlas used by the Canvas 2D scrubber.
type SpriteSheetResource struct {
	SheetURL    string `json:"sheet_url"`
	FramesCount int    `json:"frames_count"`
	LayoutRows  int    `json:"layout_rows"`
	LayoutCols  int    `json:"layout_cols"`
	FrameWidth  int    `json:"frame_width"`
	FrameHeight int    `json:"frame_height"`
	TotalSizeKB int    `json:"total_size_kb"`
}

// Look is the aggregate root representing a user's full-body outfit scan.
type Look struct {
	ID               string               `json:"id"`
	UserID           string               `json:"user_id"`
	OriginalPhotoURL string               `json:"original_photo_url"`
	DepthMapURL      string               `json:"depth_map_url,omitempty"`
	Season           Season               `json:"season"`
	MinTemp          int                  `json:"min_temp"`
	MaxTemp          int                  `json:"max_temp"`
	StyleTags        []string             `json:"style_tags"`
	Garments         []GarmentItem        `json:"garments"`
	SpriteSheet      *SpriteSheetResource `json:"sprite_sheet,omitempty"`
	CreatedAt        time.Time            `json:"created_at"`
	UpdatedAt        time.Time            `json:"updated_at"`
}

// AnalysisResult holds the structured perception output before matting and persistence.
type AnalysisResult struct {
	Season         Season
	MinTemp        int
	MaxTemp        int
	StyleTags      []string
	GarmentPrompts []GarmentCandidate
}

// GarmentCandidate holds initial LLM detection before SAM 2 / BiRefNet alpha segmentation.
type GarmentCandidate struct {
	Category     GarmentCategory
	Name         string
	Subcategory  string
	BoundingBox  BoundingBox
	ColorPalette []string
	Material     string
	Confidence   float64
}

// TryOnResult holds the 2D native-pose virtual try-on render output.
type TryOnResult struct {
	LookID             string        `json:"look_id"`
	ResultImageURL     string        `json:"result_image_url"`
	SelectedGarmentIDs []string      `json:"selected_garment_ids"`
	LatencyMS          int64         `json:"latency_ms"`
	UpdatedGarments    []GarmentItem `json:"updated_garments,omitempty"`
}

// Orbit360Result holds the completed 16-frame sprite sheet output.
type Orbit360Result struct {
	LookID      string              `json:"look_id"`
	SpriteSheet SpriteSheetResource `json:"sprite_sheet"`
	LatencyMS   int64               `json:"latency_ms"`
}
