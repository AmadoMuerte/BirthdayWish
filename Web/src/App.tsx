import { Outlet, useRouter } from '@tanstack/react-router'
import Header from './components/header/header'
import Footer from './components/footer/footer'
import styles from './app.module.css'

function App() {
  const router = useRouter()
  const currentPath = router.state.location.pathname
  const shouldShowHeader = !currentPath.startsWith('/app')

  return (
    <>
      {shouldShowHeader && <Header />}
      <main className={styles.main}>
        <Outlet />
      </main>
      <Footer />
    </>
  )
}

export default App