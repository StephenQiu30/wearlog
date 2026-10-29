// @ts-ignore
/* eslint-disable */
import request, { type RequestOptions } from '../lib/api/request'

/** 分析全身穿搭并解构单品与季节打标 POST /orbit/analyze */
export async function analyzeOrbitOutfit(
  body: API.AnalyzeOutfitRequest,
  options?: RequestOptions,
) {
  return request<API.LookResponse>('/orbit/analyze', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  })
}

/** 获取已完成的穿搭立体资产与单品列表 GET /orbit/looks/${param0} */
export async function getOrbitLook(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.getOrbitLookParams,
  options?: RequestOptions,
) {
  const { look_id: param0, ...queryParams } = params
  return request<API.LookResponse>(`/orbit/looks/${param0}`, {
    method: 'GET',
    params: { ...queryParams },
    ...(options || {}),
  })
}

/** 生成 16 帧 360° 交互式关键帧精灵图 POST /orbit/looks/${param0}/generate-360 */
export async function generateOrbit360(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.generateOrbit360Params,
  options?: RequestOptions,
) {
  const { look_id: param0, ...queryParams } = params
  return request<API.Generate360Response>(
    `/orbit/looks/${param0}/generate-360`,
    {
      method: 'POST',
      params: { ...queryParams },
      ...(options || {}),
    },
  )
}

/** 保持原生姿态的 2D 虚拟换装 POST /orbit/looks/${param0}/try-on */
export async function tryOnOrbitLook(
  // 叠加生成的Param类型 (非body参数swagger默认没有生成对象)
  params: API.tryOnOrbitLookParams,
  body: API.TryOnRequest,
  options?: RequestOptions,
) {
  const { look_id: param0, ...queryParams } = params
  return request<API.TryOnResponse>(`/orbit/looks/${param0}/try-on`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    params: { ...queryParams },
    data: body,
    ...(options || {}),
  })
}
