import styles from './footer.module.css'

export default function Footer() {

    return (
        <footer className={styles.footer}>
            <div className='container'>
                <div className={styles.footerInner}>
                    <p className={styles.copyright}>Copyright &copy; {new Date().getFullYear()}  <a href="https://github.com/AmadoMuerte/BirthdayWish">BirthdayWish</a></p>
                    <p className={styles.creator}>Created by <a href="https://github.com/AmadoMuerte">AmadoMuerte</a></p>
                </div>
            </div>
        </footer>
    )
}