import { createFileRoute } from '@tanstack/react-router'
import { ProtectedRoute } from '../components/protected_route'

export const Route = createFileRoute('/app')({
  component: () => (
    <ProtectedRoute>
      <RouteComponent />
    </ProtectedRoute>
  )
})

function RouteComponent() {
  return (
    <div className="container">
      <div>Hello "/app"!</div>
    </div>
  )
}