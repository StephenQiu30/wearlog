'use client'

import { Shirt, Footprints, Sparkles, Check, Wand2 } from 'lucide-react'
import Image from 'next/image'
import { Button } from '@/components/ui/button'
import { Spinner } from '@/components/ui/spinner'

export interface GarmentDockItem {
  id: string
  category: string
  name: string
  subcategory?: string
  matting_png_url?: string
  material?: string
}

export interface GarmentDockProps {
  garments: GarmentDockItem[]
  selectedIds: string[]
  onToggle: (id: string) => void
  onTryOn: () => void
  isTryingOn?: boolean
  className?: string
}

const categoryIcons: Record<string, React.ReactNode> = {
  outerwear: <Shirt className="h-4 w-4" />,
  top: <Shirt className="h-4 w-4" />,
  bottom: <Footprints className="h-4 w-4" />,
  shoes: <Footprints className="h-4 w-4" />,
  accessory: <Sparkles className="h-4 w-4" />,
}

export function GarmentDock({
  garments = [],
  selectedIds = [],
  onToggle,
  onTryOn,
  isTryingOn = false,
  className = '',
}: GarmentDockProps) {
  if (garments.length === 0) return null

  return (
    <div
      className={`w-full max-w-xl mx-auto p-2.5 rounded-2xl bg-background/80 backdrop-blur-xl border shadow-lg flex flex-col gap-2 ${className}`}
    >
      <div className="flex items-center justify-between px-1.5">
        <span className="text-xs font-medium text-muted-foreground">
          数字衣橱单品 ({garments.length})
        </span>
        <Button
          size="sm"
          variant="default"
          onClick={onTryOn}
          disabled={isTryingOn || selectedIds.length === 0}
          className="h-7 text-xs gap-1.5 px-3"
        >
          {isTryingOn ? (
            <>
              <Spinner className="h-3 w-3" />
              <span>生成中...</span>
            </>
          ) : (
            <>
              <Wand2 className="h-3 w-3" />
              <span>试穿选中 ({selectedIds.length})</span>
            </>
          )}
        </Button>
      </div>

      {/* Horizontal scroll items */}
      <div className="flex items-center gap-2 overflow-x-auto pb-1 pt-0.5 px-1 scrollbar-none">
        {garments.map((item) => {
          const isSelected = selectedIds.includes(item.id)
          return (
            <button
              key={item.id}
              type="button"
              onClick={() => onToggle(item.id)}
              className={`group relative flex flex-col items-center gap-1 p-1.5 rounded-xl border text-left transition-all duration-150 shrink-0 w-20 ${
                isSelected
                  ? 'border-primary ring-2 ring-primary/20 bg-primary/5'
                  : 'border-border/60 hover:border-border bg-card/50'
              }`}
            >
              {/* Thumbnail or Category Icon */}
              <div className="relative w-14 h-14 rounded-lg bg-muted/40 overflow-hidden flex items-center justify-center">
                {item.matting_png_url &&
                !item.matting_png_url.startsWith('/sample-clothes/') ? (
                  <Image
                    src={item.matting_png_url}
                    alt={item.name}
                    fill
                    sizes="56px"
                    className="object-contain p-1"
                  />
                ) : (
                  <div className="text-muted-foreground">
                    {categoryIcons[item.category] || (
                      <Shirt className="h-5 w-5" />
                    )}
                  </div>
                )}

                {/* Selected Checkmark */}
                {isSelected && (
                  <div className="absolute top-1 right-1 h-3.5 w-3.5 rounded-full bg-primary text-primary-foreground flex items-center justify-center">
                    <Check className="h-2.5 w-2.5 stroke-[3]" />
                  </div>
                )}
              </div>

              {/* Item Name */}
              <span className="text-[11px] font-medium truncate w-full text-center text-foreground">
                {item.name}
              </span>
            </button>
          )
        })}
      </div>
    </div>
  )
}
