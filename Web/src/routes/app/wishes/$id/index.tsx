import { createFileRoute } from '@tanstack/react-router'
import { WishPage } from '../../../../pages/wish/WishPage'
import { z } from 'zod'

export const Route = createFileRoute('/app/wishes/$id/')({
  params: {
    parse: (params) => ({
      id: z.number().int().parse(Number(params.id))
    }),
    stringify: ({ id }) => ({ id: `${id}` })
  },
  component: WishComponent
})

function WishComponent() {
  const { id } = Route.useParams()
  return <WishPage id={id} />
}
