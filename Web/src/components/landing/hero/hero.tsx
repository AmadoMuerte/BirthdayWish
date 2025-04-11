import { useNavigate } from '@tanstack/react-router'
import Button from '../../button/button'
import styles from './hero.module.css'

export default function Hero() {

    const navigate = useNavigate()

    return (
        <section className={styles.hero}>
            <div className="container">
                <div className={styles.heroInner}>
                    <div className={styles.heroText}>
                        <h1>Stop unwanted surprises start wishlisting!</h1>
                        <p>BirthdayWish turns gift-giving into a joy—for you and your friends. <br />Every present, perfectly picked.</p>
                    </div>
                    <div className={styles.heroForm}>
                        <div>
                            <Button handler={() => { 

                                navigate({ to: '/registration'})

                            }} text='Register' />
                            <p>Sign up & never get a bad gift again. Your future self will thank you</p>
                        </div>
                    </div>
                </div>
            </div>
        </section>
    )
}