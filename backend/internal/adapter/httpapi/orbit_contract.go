package httpapi

import (
	"time"
)

// AnalyzeOutfitRequest contains the photo URL to perceive and deconstruct into digital wardrobe items.
type AnalyzeOutfitRequest struct {
	PhotoURL string `json:"photo_url" doc:"全身穿搭照片绝对 URL"`
}

// GarmentItemDTO represents one segmented item from the outfit.
type GarmentItemDTO struct {
	ID            string     `json:"id" doc:"单品唯一标识"`
	Category      string     `json:"category" enum:"outerwear,top,bottom,shoes,accessory" doc:"衣物主类别"`
	Name          string     `json:"name" doc:"单品描述性名称"`
	Subcategory   string     `json:"subcategory" doc:"单品具体品类细分"`
	MattingPNGURL string     `json:"matting_png_url" doc:"透明抠图 PNG 访问地址"`
	BoundingBox   [4]float64 `json:"bounding_box" doc:"归一化目标检测边界框 [ymin, xmin, ymax, xmax]"`
	ColorPalette  []string   `json:"color_palette" doc:"提取的主色调 HEX 列表"`
	Material      string     `json:"material" doc:"材质面料识别标签"`
	Confidence    float64    `json:"confidence" doc:"识别置信度"`
}

// SpriteSheetDTO represents the 16-frame 360° atlas metadata.
type SpriteSheetDTO struct {
	SheetURL    string `json:"sheet_url" doc:"WebP 精灵图地址"`
	FramesCount int    `json:"frames_count" doc:"关键帧总数，默认 16"`
	LayoutRows  int    `json:"layout_rows" doc:"网格行数，默认 4"`
	LayoutCols  int    `json:"layout_cols" doc:"网格列数，默认 4"`
	FrameWidth  int    `json:"frame_width" doc:"单帧像素宽度"`
	FrameHeight int    `json:"frame_height" doc:"单帧像素高度"`
	TotalSizeKB int    `json:"total_size_kb" doc:"体积大小 (KB)"`
}

// LookResponse represents the full outfit profile aggregate root.
type LookResponse struct {
	ID               string           `json:"id" format:"uuid" doc:"Look 全局唯一标识"`
	OriginalPhotoURL string           `json:"original_photo_url" doc:"全身穿搭原图 URL"`
	DepthMapURL      string           `json:"depth_map_url,omitempty" doc:"2.5D 深度图 URL"`
	Season           string           `json:"season" enum:"spring,summer,autumn,winter,all_season" doc:"多模态季节识别标签"`
	MinTemp          int              `json:"min_temp" doc:"适穿最低温度 (摄氏度)"`
	MaxTemp          int              `json:"max_temp" doc:"适穿最高温度 (摄氏度)"`
	StyleTags        []string         `json:"style_tags" doc:"风格标签数组"`
	Garments         []GarmentItemDTO `json:"garments" doc:"拆解出的人物穿搭单品列表"`
	SpriteSheet      *SpriteSheetDTO  `json:"sprite_sheet,omitempty" doc:"360° 连续关键帧精灵图"`
	CreatedAt        time.Time        `json:"created_at" doc:"创建时间"`
	UpdatedAt        time.Time        `json:"updated_at" doc:"更新时间"`
}

// TryOnRequest holds items selected for native-pose 2D virtual try-on.
type TryOnRequest struct {
	SelectedGarmentIDs []string         `json:"selected_garment_ids" doc:"试穿单品 ID 列表"`
	Garments           []GarmentItemDTO `json:"garments" doc:"选定的试穿单品详情"`
}

// TryOnResponse returns the generated try-on composite and latency statistics.
type TryOnResponse struct {
	LookID             string   `json:"look_id" format:"uuid" doc:"关联 Look ID"`
	ResultImageURL     string   `json:"result_image_url" doc:"保持原生姿态的换装渲染图"`
	SelectedGarmentIDs []string `json:"selected_garment_ids" doc:"生效的单品 ID 列表"`
	LatencyMS          int64    `json:"latency_ms" doc:"换装渲染耗时 (毫秒)"`
}

// Generate360Response returns the generated 16-frame sprite sheet.
type Generate360Response struct {
	LookID      string         `json:"look_id" format:"uuid" doc:"关联 Look ID"`
	SpriteSheet SpriteSheetDTO `json:"sprite_sheet" doc:"16 帧 WebP 精灵图资源"`
	LatencyMS   int64          `json:"latency_ms" doc:"360 生成耗时 (毫秒)"`
}

type analyzeOutfitInput struct {
	Session string               `cookie:"then_session" hidden:"true"`
	Body    AnalyzeOutfitRequest `doc:"全身照片解析请求体"`
}

type analyzeOutfitOutput struct {
	RequestID string       `header:"X-Request-ID"`
	Body      LookResponse `doc:"穿搭深度解析与单品解构结果"`
}

type tryOnInput struct {
	Session string       `cookie:"then_session" hidden:"true"`
	LookID  string       `path:"look_id" format:"uuid" doc:"目标 Look ID"`
	Body    TryOnRequest `doc:"试穿单品请求体"`
}

type tryOnOutput struct {
	RequestID string        `header:"X-Request-ID"`
	Body      TryOnResponse `doc:"2D 原生姿态换装结果"`
}

type generate360Input struct {
	Session string `cookie:"then_session" hidden:"true"`
	LookID  string `path:"look_id" format:"uuid" doc:"目标 Look ID"`
}

type generate360Output struct {
	RequestID string              `header:"X-Request-ID"`
	Body      Generate360Response `doc:"360 关键帧精灵图生成结果"`
}

type getLookInput struct {
	Session string `cookie:"then_session" hidden:"true"`
	LookID  string `path:"look_id" format:"uuid" doc:"目标 Look ID"`
}

type getLookOutput struct {
	RequestID string       `header:"X-Request-ID"`
	Body      LookResponse `doc:"Look 详情"`
}
