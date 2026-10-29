'use client'

import { useCallback, useRef, useState } from 'react'

export interface UseOrbitGestureOptions {
  totalFrames?: number
  sensitivity?: number // Pixels per frame shift
  initialFrame?: number
}

export function useOrbitGesture({
  totalFrames = 16,
  sensitivity = 18,
  initialFrame = 0,
}: UseOrbitGestureOptions = {}) {
  const [currentFrame, setCurrentFrame] = useState<number>(initialFrame)
  const [isDragging, setIsDragging] = useState<boolean>(false)

  const startXRef = useRef<number>(0)
  const startFrameRef = useRef<number>(initialFrame)
  const isPointerDownRef = useRef<boolean>(false)

  const normalizeFrame = useCallback(
    (frame: number) => {
      return ((frame % totalFrames) + totalFrames) % totalFrames
    },
    [totalFrames],
  )

  const handlePointerDown = useCallback(
    (e: React.PointerEvent) => {
      setIsDragging(true)
      isPointerDownRef.current = true
      startXRef.current = e.clientX
      startFrameRef.current = currentFrame
      ;(e.target as HTMLElement)?.setPointerCapture?.(e.pointerId)
    },
    [currentFrame],
  )

  const handlePointerMove = useCallback(
    (e: React.PointerEvent) => {
      if (!isPointerDownRef.current) return
      const deltaX = e.clientX - startXRef.current
      const frameDelta = -Math.round(deltaX / sensitivity)
      const nextFrame = normalizeFrame(startFrameRef.current + frameDelta)
      setCurrentFrame(nextFrame)
    },
    [normalizeFrame, sensitivity],
  )

  const handlePointerUp = useCallback((e: React.PointerEvent) => {
    if (!isPointerDownRef.current) return
    isPointerDownRef.current = false
    setIsDragging(false)
    ;(e.target as HTMLElement)?.releasePointerCapture?.(e.pointerId)
  }, [])

  const handlePointerCancel = useCallback((e: React.PointerEvent) => {
    if (!isPointerDownRef.current) return
    isPointerDownRef.current = false
    setIsDragging(false)
    ;(e.target as HTMLElement)?.releasePointerCapture?.(e.pointerId)
  }, [])

  const nextFrame = useCallback(() => {
    setCurrentFrame((prev) => normalizeFrame(prev + 1))
  }, [normalizeFrame])

  const prevFrame = useCallback(() => {
    setCurrentFrame((prev) => normalizeFrame(prev - 1))
  }, [normalizeFrame])

  return {
    currentFrame,
    setFrame: setCurrentFrame,
    isDragging,
    nextFrame,
    prevFrame,
    handlers: {
      onPointerDown: handlePointerDown,
      onPointerMove: handlePointerMove,
      onPointerUp: handlePointerUp,
      onPointerCancel: handlePointerCancel,
    },
  }
}
