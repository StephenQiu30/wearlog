package orbit

import (
	"context"

	accountapp "github.com/StephenQiu30/wearlog-server/backend/internal/application/account"
)

// Authenticator verifies user session cookies and provides the current authenticated user.
type Authenticator interface {
	CurrentUser(ctx context.Context, session string) (accountapp.User, error)
}

// FashionAnalyzer perceives clothing items, seasons, temperatures and aesthetic tags from an outfit photo.
type FashionAnalyzer interface {
	Analyze(ctx context.Context, photoURL string) (*AnalysisResult, error)
}

// MattingProvider cuts out individual garment components with clean alpha channels using SAM 2 / BiRefNet.
type MattingProvider interface {
	MatGarments(ctx context.Context, photoURL string, candidates []GarmentCandidate) ([]GarmentItem, error)
}

// TryOnProvider performs 2D native-pose garment replacement via cloud diffusion API.
type TryOnProvider interface {
	TryOn(ctx context.Context, modelPhotoURL string, garments []GarmentItem) (string, error)
}

// OrbitGenerator renders or synthesizes 16 canonical azimuth keyframes and assembles a 4x4 WebP sprite sheet.
type OrbitGenerator interface {
	Generate360(ctx context.Context, photoURL string) (*SpriteSheetResource, error)
}

// Repository handles persistence of Looks and extracted garments.
type Repository interface {
	SaveLook(ctx context.Context, look *Look) error
	GetLook(ctx context.Context, lookID string) (*Look, error)
	ListLooksByUserID(ctx context.Context, userID string, limit, offset int) ([]*Look, error)
}
