import { UpdateWish } from '@/pages/update-wish'
import { createFileRoute } from '@tanstack/react-router'
import { z } from 'zod'

export const Route = createFileRoute('/app/wishes/$id/update')({
  params: {
    parse: (params) => ({
      id: z.number().int().parse(Number(params.id))
    }),
    stringify: ({ id }) => ({ id: `${id}` })
  },
  component: UpdateWishComponent
})

function UpdateWishComponent() {
  const { id } = Route.useParams()
  return (
    <div>
      <UpdateWish id={id} />
    </div>
  )
}
