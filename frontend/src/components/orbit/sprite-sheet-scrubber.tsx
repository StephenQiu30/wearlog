'use client'

import { useEffect, useRef, useState } from 'react'
import { Rotate3d } from 'lucide-react'

export interface SpriteSheetScrubberProps {
  sheetURL: string
  framesCount?: number
  layoutCols?: number
  layoutRows?: number
  frameWidth?: number
  frameHeight?: number
  currentFrame: number
  isDragging?: boolean
  handlers?: {
    onPointerDown: (e: React.PointerEvent) => void
    onPointerMove: (e: React.PointerEvent) => void
    onPointerUp: (e: React.PointerEvent) => void
    onPointerCancel: (e: React.PointerEvent) => void
  }
  className?: string
}

export function SpriteSheetScrubber({
  sheetURL,
  layoutCols = 4,
  frameWidth = 512,
  frameHeight = 683,
  currentFrame,
  isDragging = false,
  handlers,
  className = '',
}: SpriteSheetScrubberProps) {
  const canvasRef = useRef<HTMLCanvasElement | null>(null)
  const imageRef = useRef<HTMLImageElement | null>(null)
  const [imageLoaded, setImageLoaded] = useState<boolean>(false)
  const [loadError, setLoadError] = useState<boolean>(false)

  // Preload sprite sheet image
  useEffect(() => {
    let active = true
    const img = new Image()
    img.crossOrigin = 'anonymous'
    img.src = sheetURL
    img.onload = () => {
      if (!active) return
      imageRef.current = img
      setImageLoaded(true)
      setLoadError(false)
    }
    img.onerror = () => {
      if (!active) return
      setLoadError(true)
    }

    return () => {
      active = false
    }
  }, [sheetURL])

  // Canvas 2D Blit on frame or size update
  useEffect(() => {
    const canvas = canvasRef.current
    const img = imageRef.current
    if (!canvas || !img || !imageLoaded) return

    const ctx = canvas.getContext('2d', { alpha: true })
    if (!ctx) return

    const col = currentFrame % layoutCols
    const row = Math.floor(currentFrame / layoutCols)
    const sx = col * frameWidth
    const sy = row * frameHeight

    const render = () => {
      ctx.clearRect(0, 0, canvas.width, canvas.height)
      ctx.drawImage(
        img,
        sx,
        sy,
        frameWidth,
        frameHeight,
        0,
        0,
        canvas.width,
        canvas.height,
      )
    }

    const animId = requestAnimationFrame(render)
    return () => cancelAnimationFrame(animId)
  }, [currentFrame, imageLoaded, layoutCols, frameWidth, frameHeight])

  return (
    <div
      className={`relative w-full h-full select-none touch-none cursor-grab active:cursor-grabbing ${className}`}
      {...handlers}
    >
      <canvas
        ref={canvasRef}
        width={frameWidth}
        height={frameHeight}
        className="w-full h-full object-contain pointer-events-none"
      />

      {/* Frame indicator badge */}
      <div className="absolute bottom-3 right-3 px-2 py-0.5 rounded-full bg-black/40 backdrop-blur-md text-[10px] text-white/90 font-mono flex items-center gap-1 pointer-events-none transition-opacity duration-200">
        <Rotate3d className="h-3 w-3" />
        <span>{Math.round((currentFrame / 16) * 360)}°</span>
      </div>

      {/* Touch drag prompt */}
      {!isDragging && (
        <div className="absolute top-3 left-1/2 -translate-x-1/2 px-2.5 py-1 rounded-full bg-black/30 backdrop-blur-sm text-xs text-white/80 pointer-events-none flex items-center gap-1.5 animate-pulse">
          <Rotate3d className="h-3.5 w-3.5" />
          <span>左右拖动 360° 查看</span>
        </div>
      )}

      {loadError && (
        <div className="absolute inset-0 flex items-center justify-center bg-muted/40 text-xs text-muted-foreground">
          关键帧精灵图加载中...
        </div>
      )}
    </div>
  )
}
