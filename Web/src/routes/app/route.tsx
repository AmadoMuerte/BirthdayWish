import { createFileRoute, Outlet } from '@tanstack/react-router'
import { ProtectedRoute } from '../../components/protected_route'
import { MantineProvider } from '@mantine/core'
import '@mantine/core/styles.css'
import { AppWrapper } from '../../components/app_wrapper/app_wrapper'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'

export const Route = createFileRoute('/app')({
  component: () => (
    <ProtectedRoute>
      <RouteComponent />
    </ProtectedRoute>
  )
})

function RouteComponent() {
  const queryClient = new QueryClient()
  return (
    <QueryClientProvider client={queryClient}>
      <MantineProvider>
        <AppWrapper content={<Outlet />} />
      </MantineProvider>
    </QueryClientProvider>
  )
}
