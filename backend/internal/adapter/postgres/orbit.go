package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/StephenQiu30/wearlog-server/backend/internal/application/orbit"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrbitRepository struct {
	database *gorm.DB
}

func NewOrbitRepository(database *gorm.DB) *OrbitRepository {
	return &OrbitRepository{database: database}
}

type orbitLookRecord struct {
	ID               string    `gorm:"column:id;type:uuid;primaryKey"`
	UserID           string    `gorm:"column:user_id;type:uuid;index;not null"`
	OriginalPhotoURL string    `gorm:"column:original_photo_url;type:text;not null"`
	DepthMapURL      string    `gorm:"column:depth_map_url;type:text;not null;default:''"`
	Season           string    `gorm:"column:season;type:varchar(32);not null"`
	MinTemp          int       `gorm:"column:min_temp;not null"`
	MaxTemp          int       `gorm:"column:max_temp;not null"`
	StyleTagsJSON    string    `gorm:"column:style_tags;type:text;not null;default:'[]'"`
	GarmentsJSON     string    `gorm:"column:garments;type:text;not null;default:'[]'"`
	SpriteSheetJSON  string    `gorm:"column:sprite_sheet;type:text;not null;default:''"`
	CreatedAt        time.Time `gorm:"column:created_at;type:timestamptz;not null"`
	UpdatedAt        time.Time `gorm:"column:updated_at;type:timestamptz;not null"`
}

func (orbitLookRecord) TableName() string {
	return "orbit_looks"
}

func (r *OrbitRepository) SaveLook(ctx context.Context, look *orbit.Look) error {
	if r.database == nil {
		return errors.New("database not available")
	}
	if look == nil {
		return errors.New("cannot save nil look")
	}

	styleTagsBytes, err := json.Marshal(look.StyleTags)
	if err != nil {
		return fmt.Errorf("marshal style tags: %w", err)
	}
	garmentsBytes, err := json.Marshal(look.Garments)
	if err != nil {
		return fmt.Errorf("marshal garments: %w", err)
	}
	var spriteSheetStr string
	if look.SpriteSheet != nil {
		b, err := json.Marshal(look.SpriteSheet)
		if err != nil {
			return fmt.Errorf("marshal sprite sheet: %w", err)
		}
		spriteSheetStr = string(b)
	}

	record := orbitLookRecord{
		ID:               look.ID,
		UserID:           look.UserID,
		OriginalPhotoURL: look.OriginalPhotoURL,
		DepthMapURL:      look.DepthMapURL,
		Season:           string(look.Season),
		MinTemp:          look.MinTemp,
		MaxTemp:          look.MaxTemp,
		StyleTagsJSON:    string(styleTagsBytes),
		GarmentsJSON:     string(garmentsBytes),
		SpriteSheetJSON:  spriteSheetStr,
		CreatedAt:        look.CreatedAt,
		UpdatedAt:        look.UpdatedAt,
	}

	return r.database.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"original_photo_url", "depth_map_url", "season", "min_temp", "max_temp",
			"style_tags", "garments", "sprite_sheet", "updated_at",
		}),
	}).Create(&record).Error
}

func (r *OrbitRepository) GetLook(ctx context.Context, lookID string) (*orbit.Look, error) {
	if r.database == nil {
		return nil, errors.New("database not available")
	}
	var record orbitLookRecord
	err := r.database.WithContext(ctx).First(&record, "id = ?", lookID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return recordToLook(&record)
}

func (r *OrbitRepository) ListLooksByUserID(ctx context.Context, userID string, limit, offset int) ([]*orbit.Look, error) {
	if r.database == nil {
		return nil, errors.New("database not available")
	}
	var records []orbitLookRecord
	query := r.database.WithContext(ctx).Where("user_id = ?", userID).Order("created_at desc")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	if err := query.Find(&records).Error; err != nil {
		return nil, err
	}

	looks := make([]*orbit.Look, 0, len(records))
	for i := range records {
		look, err := recordToLook(&records[i])
		if err != nil {
			return nil, err
		}
		looks = append(looks, look)
	}
	return looks, nil
}

func recordToLook(record *orbitLookRecord) (*orbit.Look, error) {
	var styleTags []string
	if record.StyleTagsJSON != "" {
		if err := json.Unmarshal([]byte(record.StyleTagsJSON), &styleTags); err != nil {
			return nil, fmt.Errorf("unmarshal style tags: %w", err)
		}
	}

	var garments []orbit.GarmentItem
	if record.GarmentsJSON != "" {
		if err := json.Unmarshal([]byte(record.GarmentsJSON), &garments); err != nil {
			return nil, fmt.Errorf("unmarshal garments: %w", err)
		}
	}

	var spriteSheet *orbit.SpriteSheetResource
	if record.SpriteSheetJSON != "" {
		var s orbit.SpriteSheetResource
		if err := json.Unmarshal([]byte(record.SpriteSheetJSON), &s); err != nil {
			return nil, fmt.Errorf("unmarshal sprite sheet: %w", err)
		}
		spriteSheet = &s
	}

	return &orbit.Look{
		ID:               record.ID,
		UserID:           record.UserID,
		OriginalPhotoURL: record.OriginalPhotoURL,
		DepthMapURL:      record.DepthMapURL,
		Season:           orbit.Season(record.Season),
		MinTemp:          record.MinTemp,
		MaxTemp:          record.MaxTemp,
		StyleTags:        styleTags,
		Garments:         garments,
		SpriteSheet:      spriteSheet,
		CreatedAt:        record.CreatedAt,
		UpdatedAt:        record.UpdatedAt,
	}, nil
}
