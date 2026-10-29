'use client'

import { useEffect, useRef, useState } from 'react'
import Image from 'next/image'

export interface ParallaxCanvasProps {
  photoURL: string
  depthMapURL?: string
  className?: string
}

export function ParallaxCanvas({
  photoURL,
  depthMapURL,
  className = '',
}: ParallaxCanvasProps) {
  const containerRef = useRef<HTMLDivElement | null>(null)
  const [offset, setOffset] = useState<{ x: number; y: number }>({ x: 0, y: 0 })
  const [hasDepth, setHasDepth] = useState<boolean>(Boolean(depthMapURL))
  const targetOffsetRef = useRef<{ x: number; y: number }>({ x: 0, y: 0 })

  // Smooth lerp loop for tactile inertia
  useEffect(() => {
    let animId: number
    const loop = () => {
      setOffset((current) => {
        const target = targetOffsetRef.current
        const nextX = current.x + (target.x - current.x) * 0.1
        const nextY = current.y + (target.y - current.y) * 0.1
        return { x: nextX, y: nextY }
      })
      animId = requestAnimationFrame(loop)
    }
    animId = requestAnimationFrame(loop)
    return () => cancelAnimationFrame(animId)
  }, [])

  // Optional: dynamically load onnxruntime-web in the background only after mount
  useEffect(() => {
    if (depthMapURL) return
    let active = true

    // Dynamic import to satisfy Vercel bundle-size and lazy-loading rules
    import('onnxruntime-web')
      .then((ort) => {
        if (!active) return
        // Configure WebAssembly paths from public directory if initialized
        ort.env.wasm.numThreads = 1
        setHasDepth(true)
      })
      .catch(() => {
        // Graceful fallback to 2.5D perspective tilt
      })

    return () => {
      active = false
    }
  }, [depthMapURL])

  const handlePointerMove = (e: React.PointerEvent<HTMLDivElement>) => {
    const rect = containerRef.current?.getBoundingClientRect()
    if (!rect) return
    const clientX = e.clientX - rect.left
    const clientY = e.clientY - rect.top
    const normX = (clientX / rect.width - 0.5) * 2 // -1 to 1
    const normY = (clientY / rect.height - 0.5) * 2 // -1 to 1

    targetOffsetRef.current = {
      x: normX * 12, // max 12px shift
      y: normY * 8, // max 8px shift
    }
  }

  const handlePointerLeave = () => {
    targetOffsetRef.current = { x: 0, y: 0 }
  }

  return (
    <div
      ref={containerRef}
      onPointerMove={handlePointerMove}
      onPointerLeave={handlePointerLeave}
      className={`relative w-full h-full overflow-hidden select-none cursor-crosshair ${className}`}
    >
      {/* 2.5D Parallax layered image */}
      <div
        className="absolute inset-0 transition-transform duration-75 ease-out will-change-transform"
        style={{
          transform: `translate3d(${offset.x}px, ${offset.y}px, 0px) scale(1.05)`,
        }}
      >
        <Image
          src={photoURL}
          alt="人物全身立体穿搭"
          fill
          sizes="(max-width: 768px) 100vw, 480px"
          priority
          className="object-contain"
        />
      </div>

      {/* Tactile hint badge */}
      <div className="absolute top-3 left-1/2 -translate-x-1/2 px-2.5 py-1 rounded-full bg-black/30 backdrop-blur-sm text-xs text-white/80 pointer-events-none flex items-center gap-1.5">
        <span>
          {hasDepth ? '2.5D 深度视差 · 移动鼠标探索' : '立体视差模式'}
        </span>
      </div>
    </div>
  )
}
