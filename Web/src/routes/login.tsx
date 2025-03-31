import { createFileRoute } from '@tanstack/react-router'
import LoginForm from '../components/login_form/login_form'

export const Route = createFileRoute('/login')({
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <div className='loginPage'>
      <LoginForm />
    </div>
  )
}
