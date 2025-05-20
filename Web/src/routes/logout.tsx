import { createFileRoute } from '@tanstack/react-router'
import AuthForm from '../components/auth_form/auth_form'
import { removeToken } from '../api/token'

export const Route = createFileRoute('/logout')({
  component: RouteComponent,
})

function RouteComponent() {
  removeToken();

  return (
    <div className='loginPage'>
      <AuthForm isRegistration={false} />
    </div>
  )
}
