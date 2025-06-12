import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/app/friends')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/app/friends"!</div>
}
