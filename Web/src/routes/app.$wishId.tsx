import { createFileRoute } from '@tanstack/react-router'
import { z } from 'zod'

export const Route = createFileRoute('/app/$wishId')({
  params: {
    parse: (params) => ({
      wishId: z.number().int().parse(Number(params.wishId)),
    }),
    stringify: ({ wishId }) => ({ wishId: `${wishId}` }),
  },
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/app/$wishId"!</div>
}
