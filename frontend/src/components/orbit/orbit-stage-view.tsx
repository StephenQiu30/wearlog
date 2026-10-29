'use client'

import { useState } from 'react'
import { Rotate3d, Layers, Sparkles, Upload } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Spinner } from '@/components/ui/spinner'
import { useOrbit } from '@/hooks/orbit/use-orbit'
import { useOrbitGesture } from '@/hooks/orbit/use-orbit-gesture'
import { ParallaxCanvas } from './parallax-canvas'
import { SpriteSheetScrubber } from './sprite-sheet-scrubber'
import { GarmentDock } from './garment-dock'
import { OutfitAttributePill } from './outfit-attribute-pill'

export interface OrbitStageViewProps {
  initialPhotoURL?: string
  className?: string
}

const DEFAULT_SAMPLE_PHOTO =
  'https://images.unsplash.com/photo-1515886657613-9f3515b0c78f?w=800&auto=format&fit=crop&q=80'

export function OrbitStageView({
  initialPhotoURL = DEFAULT_SAMPLE_PHOTO,
  className = '',
}: OrbitStageViewProps) {
  const [photoURL, setPhotoURL] = useState<string>(initialPhotoURL)
  const {
    currentLook,
    mode,
    setMode,
    selectedGarmentIds,
    toggleGarment,
    isAnalyzing,
    isTryingOn,
    isGenerating360,
    errorMessage,
    analyze,
    tryOn,
    request360,
  } = useOrbit()

  const { currentFrame, isDragging, handlers } = useOrbitGesture({
    totalFrames: currentLook?.sprite_sheet?.frames_count || 16,
  })

  const hasSpriteSheet = Boolean(currentLook?.sprite_sheet?.sheet_url)

  const handleStartAnalysis = async () => {
    if (!photoURL) return
    await analyze(photoURL)
  }

  return (
    <div
      className={`relative flex flex-col items-center w-full min-h-[calc(100vh-4rem)] p-4 md:p-6 ${className}`}
    >
      {/* Top Header Controls */}
      <div className="w-full max-w-xl flex flex-col gap-2 mb-3">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-lg font-semibold tracking-tight">
              立体数字衣橱
            </h1>
            <p className="text-xs text-muted-foreground">
              拖动查看 360° 形象与一键拆解单品
            </p>
          </div>

          <div className="flex items-center gap-2">
            {!currentLook ? (
              <Button
                size="sm"
                onClick={handleStartAnalysis}
                disabled={isAnalyzing}
                className="gap-1.5 h-8 text-xs font-medium"
              >
                {isAnalyzing ? (
                  <>
                    <Spinner className="h-3 w-3" />
                    <span>AI 解构中...</span>
                  </>
                ) : (
                  <>
                    <Sparkles className="h-3.5 w-3.5" />
                    <span>开始解构</span>
                  </>
                )}
              </Button>
            ) : hasSpriteSheet ? (
              <div className="flex items-center rounded-lg border p-0.5 bg-muted/40">
                <Button
                  size="sm"
                  variant={mode === 'depth-parallax' ? 'secondary' : 'ghost'}
                  onClick={() => setMode('depth-parallax')}
                  className="h-7 text-xs px-2.5 gap-1"
                >
                  <Layers className="h-3.5 w-3.5" />
                  <span>2.5D</span>
                </Button>
                <Button
                  size="sm"
                  variant={mode === 'orbit-360' ? 'secondary' : 'ghost'}
                  onClick={() => setMode('orbit-360')}
                  className="h-7 text-xs px-2.5 gap-1"
                >
                  <Rotate3d className="h-3.5 w-3.5" />
                  <span>360°</span>
                </Button>
              </div>
            ) : (
              <Button
                size="sm"
                variant="outline"
                onClick={request360}
                disabled={isGenerating360}
                className="gap-1.5 h-8 text-xs font-medium"
              >
                {isGenerating360 ? (
                  <>
                    <Spinner className="h-3 w-3" />
                    <span>生成 360°...</span>
                  </>
                ) : (
                  <>
                    <Rotate3d className="h-3.5 w-3.5" />
                    <span>生成 360° 环绕</span>
                  </>
                )}
              </Button>
            )}
          </div>
        </div>

        {/* Outfit Attributes & Tags */}
        {currentLook && (
          <OutfitAttributePill
            season={currentLook.season}
            minTemp={currentLook.min_temp}
            maxTemp={currentLook.max_temp}
            styleTags={currentLook.style_tags}
          />
        )}
      </div>

      {/* Main Stage Viewport (Strict aspect-[3/4] preventing CLS) */}
      <div className="relative w-full max-w-sm sm:max-w-md aspect-[3/4] rounded-3xl overflow-hidden border bg-gradient-to-b from-muted/20 via-background to-muted/30 shadow-2xl flex items-center justify-center">
        {mode === 'orbit-360' && hasSpriteSheet && currentLook?.sprite_sheet ? (
          <SpriteSheetScrubber
            sheetURL={currentLook.sprite_sheet.sheet_url}
            framesCount={currentLook.sprite_sheet.frames_count}
            layoutCols={currentLook.sprite_sheet.layout_cols}
            layoutRows={currentLook.sprite_sheet.layout_rows}
            frameWidth={currentLook.sprite_sheet.frame_width}
            frameHeight={currentLook.sprite_sheet.frame_height}
            currentFrame={currentFrame}
            isDragging={isDragging}
            handlers={handlers}
          />
        ) : (
          <ParallaxCanvas
            photoURL={currentLook?.original_photo_url || photoURL}
            depthMapURL={currentLook?.depth_map_url}
          />
        )}

        {/* Contact Shadow at Baseline 88% */}
        <div
          className="absolute left-1/2 -translate-x-1/2 w-48 h-6 rounded-[100%] bg-black/15 blur-md pointer-events-none"
          style={{ top: '88%' }}
        />

        {/* Photo URL Switcher Modal / Overlay if not yet analyzed */}
        {!currentLook && (
          <div className="absolute inset-x-3 bottom-3 p-3 rounded-2xl bg-background/80 backdrop-blur-md border flex items-center justify-between text-xs">
            <span className="text-muted-foreground truncate mr-2">
              示例照片已就绪
            </span>
            <Button
              size="sm"
              variant="outline"
              onClick={() => {
                const next =
                  photoURL === DEFAULT_SAMPLE_PHOTO
                    ? 'https://images.unsplash.com/photo-1509631179647-0177331693ae?w=800&auto=format&fit=crop&q=80'
                    : DEFAULT_SAMPLE_PHOTO
                setPhotoURL(next)
              }}
              className="h-7 text-xs gap-1"
            >
              <Upload className="h-3 w-3" />
              <span>切换示例</span>
            </Button>
          </div>
        )}
      </div>

      {/* Error Message Notice */}
      {errorMessage && (
        <div className="w-full max-w-xl mt-3 p-2.5 rounded-xl border border-destructive/20 bg-destructive/10 text-destructive text-xs text-center">
          {errorMessage}
        </div>
      )}

      {/* Bottom Floating Garment Dock */}
      {currentLook && currentLook.garments && (
        <div className="w-full max-w-xl mt-4">
          <GarmentDock
            garments={currentLook.garments}
            selectedIds={selectedGarmentIds}
            onToggle={toggleGarment}
            onTryOn={tryOn}
            isTryingOn={isTryingOn}
          />
        </div>
      )}
    </div>
  )
}
