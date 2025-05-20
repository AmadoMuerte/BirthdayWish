import { createFileRoute, Outlet } from '@tanstack/react-router'
import { ProtectedRoute } from '../components/protected_route'
import { MantineProvider } from '@mantine/core'
import '@mantine/core/styles.css';
import { AppWrapper } from '../components/app_wrapper/app_wrapper';

export const Route = createFileRoute('/app')({
  component: () => (
    <ProtectedRoute>
      <RouteComponent />
    </ProtectedRoute>
  )
})

function RouteComponent() {
  return (
    <MantineProvider>
      <AppWrapper content={<Outlet />} />
    </MantineProvider>
  )
}