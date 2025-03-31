import Button from '../button/button';
import styles from './login_form.module.css'

export default function LoginForm() {
    return (
        <div className={styles.form}>
            <h2>Login</h2>
            <form className={styles.formInner}>
                <input type="text" placeholder="Username" required />
                <input type="password" placeholder="Password" required />
                <Button handler={() => { }} text="Login" />
                <div className={styles.links}>
                    <a href="#" className={styles.link}>Forgot password?</a>
                    <a href="#" className={styles.link}>Create account</a>
                </div>
            </form>
        </div>
    );
}