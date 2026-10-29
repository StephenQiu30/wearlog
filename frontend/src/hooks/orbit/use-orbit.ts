'use client'

import { useCallback, useState } from 'react'
import {
  analyzeOrbitOutfit,
  generateOrbit360,
  tryOnOrbitLook,
} from '@/api/orbit'

export type OrbitStageMode = 'depth-parallax' | 'orbit-360'

export function useOrbit(initialLook: API.LookResponse | null = null) {
  const [currentLook, setCurrentLook] = useState<API.LookResponse | null>(
    initialLook,
  )
  const [mode, setMode] = useState<OrbitStageMode>('depth-parallax')
  const [selectedGarmentIds, setSelectedGarmentIds] = useState<string[]>([])
  const [isAnalyzing, setIsAnalyzing] = useState<boolean>(false)
  const [isTryingOn, setIsTryingOn] = useState<boolean>(false)
  const [isGenerating360, setIsGenerating360] = useState<boolean>(false)
  const [errorMessage, setErrorMessage] = useState<string | null>(null)

  const analyze = useCallback(async (photoURL: string) => {
    setIsAnalyzing(true)
    setErrorMessage(null)
    try {
      const response = await analyzeOrbitOutfit({ photo_url: photoURL })
      setCurrentLook(response)
      if (response.garments && Array.isArray(response.garments)) {
        setSelectedGarmentIds(
          response.garments.map((g: API.GarmentItemDTO) => g.id),
        )
      }
      return response
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : '穿搭解析失败'
      setErrorMessage(msg)
      throw err
    } finally {
      setIsAnalyzing(false)
    }
  }, [])

  const toggleGarment = useCallback((garmentId: string) => {
    setSelectedGarmentIds((prev) =>
      prev.includes(garmentId)
        ? prev.filter((id) => id !== garmentId)
        : [...prev, garmentId],
    )
  }, [])

  const tryOn = useCallback(async () => {
    if (!currentLook) return
    setIsTryingOn(true)
    setErrorMessage(null)
    try {
      const garmentsToWear = (currentLook.garments || []).filter(
        (g: API.GarmentItemDTO) => selectedGarmentIds.includes(g.id),
      )
      const res = await tryOnOrbitLook(
        { look_id: currentLook.id },
        {
          selected_garment_ids: selectedGarmentIds,
          garments: garmentsToWear,
        },
      )
      setCurrentLook((prev) =>
        prev
          ? {
              ...prev,
              original_photo_url: res.result_image_url,
            }
          : null,
      )
      return res
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : '换装失败'
      setErrorMessage(msg)
      throw err
    } finally {
      setIsTryingOn(false)
    }
  }, [currentLook, selectedGarmentIds])

  const request360 = useCallback(async () => {
    if (!currentLook) return
    setIsGenerating360(true)
    setErrorMessage(null)
    try {
      const res = await generateOrbit360({ look_id: currentLook.id })
      setCurrentLook((prev) =>
        prev
          ? {
              ...prev,
              sprite_sheet: res.sprite_sheet,
            }
          : null,
      )
      setMode('orbit-360')
      return res
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : '360° 生成失败'
      setErrorMessage(msg)
      throw err
    } finally {
      setIsGenerating360(false)
    }
  }, [currentLook])

  return {
    currentLook,
    setCurrentLook,
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
  }
}
