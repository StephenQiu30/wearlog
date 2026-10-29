'use client'

import { Thermometer, Sparkles } from 'lucide-react'
import { Badge } from '@/components/ui/badge'

export interface OutfitAttributePillProps {
  season?: string
  minTemp?: number
  maxTemp?: number
  styleTags?: string[]
}

const seasonLabels: Record<string, string> = {
  spring: '春季',
  summer: '夏季',
  autumn: '秋季',
  winter: '冬季',
  all_season: '四季适穿',
}

export function OutfitAttributePill({
  season,
  minTemp,
  maxTemp,
  styleTags = [],
}: OutfitAttributePillProps) {
  const seasonText = season ? seasonLabels[season] || season : null
  const hasTemp = minTemp !== undefined && maxTemp !== undefined

  return (
    <div className="flex flex-wrap items-center gap-1.5 py-1">
      {seasonText && (
        <Badge variant="outline" className="text-xs font-medium px-2 py-0.5">
          {seasonText}
        </Badge>
      )}

      {hasTemp && (
        <Badge
          variant="secondary"
          className="text-xs font-normal gap-1 px-2 py-0.5"
        >
          <Thermometer className="h-3 w-3 text-muted-foreground" />
          <span>
            {minTemp}°C ~ {maxTemp}°C
          </span>
        </Badge>
      )}

      {styleTags.map((tag) => (
        <Badge
          key={tag}
          variant="outline"
          className="text-xs text-muted-foreground gap-1 px-2 py-0.5"
        >
          <Sparkles className="h-2.5 w-2.5" />
          <span>{tag}</span>
        </Badge>
      ))}
    </div>
  )
}
