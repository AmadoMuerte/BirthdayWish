import { Outlet } from '@tanstack/react-router'
import Header from './components/header/header'
import Footer from './components/footer/footer'
import styles from './app.module.css'

function App() {

  return (
    <>
      <Header />
      <main className={styles.main}>
        <Outlet />
      </main>
      <Footer />
    </>
  )
}

export default App
