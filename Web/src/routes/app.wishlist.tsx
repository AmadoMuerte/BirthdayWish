import { createFileRoute } from '@tanstack/react-router'
import { Wishlist } from '../components/wishlist/wishlist'

export const Route = createFileRoute('/app/wishlist')({
  component: RouteComponent,
})

function RouteComponent() {
  return <Wishlist />
}
