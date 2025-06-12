import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { Wish } from './types'
import {
  createWish,
  CreateWishParams,
  deleteWish,
  DeleteWishParams,
  getAllWishes,
  GetAllWishesParams,
  getWish,
  GetWishParams,
  updateWish,
  UpdateWishParams
} from './api'

const WISH_QUERY_KEY = 'wishes'

export function useGetAllWishesApi(params: GetAllWishesParams) {
  return useQuery({
    queryKey: [WISH_QUERY_KEY, params.query],
    queryFn: () => getAllWishes(params).then((res) => res.data),
    staleTime: 5 * 60 * 1000
  })
}

export function useGetWishApi(params: GetWishParams) {
  return useQuery({
    queryKey: [WISH_QUERY_KEY, params.id],
    queryFn: () => getWish(params).then((res) => res.data)
  })
}

export function useCreateWishApi(params: CreateWishParams) {
  const queryClient = useQueryClient()

  return useMutation<Wish, Error, CreateWishParams>({
    mutationKey: [WISH_QUERY_KEY],
    mutationFn: () => createWish(params).then((res) => res.data),
    onSuccess: (newWish) => {
      queryClient.setQueryData<Wish[]>([WISH_QUERY_KEY], (oldWishes = []) => [
        ...oldWishes,
        newWish
      ])
    }
  })
}

export function useUpdateWishApi(params: UpdateWishParams) {
  const queryClient = useQueryClient()

  return useMutation<Wish, Error, UpdateWishParams>({
    mutationKey: [WISH_QUERY_KEY],
    mutationFn: () => updateWish(params).then((res) => res.data),
    onSuccess: (updatedWish, params) => {
      queryClient.setQueryData<Wish>([WISH_QUERY_KEY, params.id], updatedWish)
      queryClient.setQueryData<Wish[]>([WISH_QUERY_KEY], (oldWishes = []) =>
        oldWishes.map((wish) => (wish.id === params.id ? updatedWish : wish))
      )
    }
  })
}

export function useDeleteWishApi() {
  const queryClient = useQueryClient()

  return useMutation<Wish, Error, DeleteWishParams>({
    mutationKey: [WISH_QUERY_KEY],
    mutationFn: (params) => deleteWish(params).then((res) => res.data),
    onSuccess: (_, params) => {
      queryClient.setQueryData<Wish[]>([WISH_QUERY_KEY], (oldWishes = []) =>
        oldWishes.filter((wish) => wish.id !== params.id)
      )
    }
  })
}
