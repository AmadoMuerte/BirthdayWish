import { createFileRoute } from '@tanstack/react-router'
import CreateForm from '../components/create_form/create_form'

export const Route = createFileRoute('/app/create')({
  component: RouteComponent,
})

function RouteComponent() {
  return <CreateForm />
}
