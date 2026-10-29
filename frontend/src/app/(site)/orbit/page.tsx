import type { Metadata } from 'next'
import { OrbitStageView } from '@/components/orbit/orbit-stage-view'

export const metadata: Metadata = {
  title: '立体数字衣橱 · 360° · 衣序',
  description:
    '将穿搭照片转化为 360° 可拖动查看的立体形象，拆解单品与智能换装。',
}

export default function OrbitPage() {
  return <OrbitStageView />
}
