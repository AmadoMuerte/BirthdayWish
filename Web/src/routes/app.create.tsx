import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/app/create')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/app/create_item"!</div>
}
