import { useGetWishApi } from '@/entities/wish/model/hooks'
import { WishForm } from '@/feautures/wish-form'
import { useRouter } from '@tanstack/react-router'
import React from 'react'

interface UpdateProps {
    id: number
  }

export const UpdateWish: React.FC<UpdateProps> = ({ id }) => {
  const router = useRouter()
  const { data, isLoading, error } = useGetWishApi({ id })

  function handleBack() {
    router.history.back()
  }

  if (isLoading) {
    return <div>Loading...</div>
  }

  if (error) {
    return <div>Error: {error.message}</div>
  }

  return (
    <div>
      <WishForm initialValues={data} />
    </div>
  )
}
