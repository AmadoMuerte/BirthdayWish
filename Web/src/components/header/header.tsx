import Logo from '../logo/logo'
import Nav from '../nav/nav'
import styles from './header.module.css'

export default function Header() {

    return (
        <header className={styles.header}>
            <div className='container'>
                <div className={styles.headerInner}>
                    <Logo />
                    <Nav />
                </div>
            </div>
        </header>
    )
}