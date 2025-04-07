import { Navigate, useRouter } from '@tanstack/react-router'
import { getTokenInfo } from '../api/token'


export function ProtectedRoute({ children }: { children: React.ReactNode }) {
    const router = useRouter()
    const tokenInfo = getTokenInfo()

    if (!tokenInfo.isValid) {
        return <Navigate to="/login" search={{ redirect: router.state.location.pathname }} replace />
    }

    return <>{children}</>
}