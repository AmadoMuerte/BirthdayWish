import { createFileRoute } from '@tanstack/react-router'
import AuthForm from '../components/auth_form/auth_form'

export const Route = createFileRoute('/registration')({
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <div className='loginPage'>
      <AuthForm isRegistration={true} />
    </div>
  )
}
