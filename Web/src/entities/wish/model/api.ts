import { AxiosRequestConfig } from 'axios'
import { api } from '@/shared/api/instance'
import { CreateWishData, UpdateWishData, Wish } from './types'

export interface ApiParamsBase {
  config?: Omit<AxiosRequestConfig, 'params'>
}

export interface GetAllWishesParams {
  sort?: 'asc' | 'desc'
  filter?: string
  page?: number
  limit?: number
}

export interface GetAllWishesParams extends ApiParamsBase {
  query?: GetAllWishesParams
}

export interface CreateWishParams extends ApiParamsBase {
  body: CreateWishData
}

export interface GetWishParams extends ApiParamsBase {
  id: number
}

export interface UpdateWishParams extends GetWishParams {
  id: number
  body: UpdateWishData
}

export interface DeleteWishParams extends GetWishParams {
  id: number
}

export const getAllWishes = ({ config, query }: GetAllWishesParams) =>
  api.get<Wish[]>('/api/wish', { ...config, params: query })

export const createWish = ({ body, config }: CreateWishParams) =>
  api.post<Wish>('api/wish', body, config)

export const getWish = ({ id, config }: GetWishParams) =>
  api.get<Wish>(`/api/wish/${id}`, config)

export const updateWish = ({ id, body, config }: UpdateWishParams) =>
  api.put<Wish>(`api/wish/${id}`, body, config)

export const deleteWish = ({ id, config }: DeleteWishParams) =>
  api.delete<Wish>(`api/wish/${id}`, config)
