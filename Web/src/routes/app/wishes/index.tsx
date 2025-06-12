import { createFileRoute } from '@tanstack/react-router'
import { Wishlist } from '@/pages/wishlist/Wishlist'

export const Route = createFileRoute('/app/wishes/')({
  component: RouteComponent,
})

function RouteComponent() {


  return (
    <Wishlist />
  )
}
