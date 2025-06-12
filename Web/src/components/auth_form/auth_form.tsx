import { Link } from '@tanstack/react-router';
import Button from '../button/button';
import styles from './auth_form.module.css';
import { useState } from 'react';
import { AuthResponse, Form, loginUser, registerUser } from '../../shared/api/auth'
import validateEmail, { validatePassword, validateUsername } from '../../shared/lib/validator';

interface Props {
    isRegistration: boolean;
}

export default function AuthForm(props: Props) {
    const { isRegistration } = props;

    const [form, setForm] = useState<Form>({
        username: '',
        password: '',
        email: isRegistration ? '' : undefined,
    });

    function handleChange(event: React.ChangeEvent<HTMLInputElement>) {
        const { name, value } = event.target;
        setForm((prevForm) => ({
            ...prevForm,
            [name]: value,
        }));
    }

    async function handleClick(event: React.MouseEvent<HTMLButtonElement>) {
        event.preventDefault();

        if (!validatePassword(form.password)) {
            // CHANGE ALERT ON CASTOM
            alert('Password is not valid.');
            return
        }

        if (!validateUsername(form.username)) {
            // CHANGE ALERT ON CASTOM
            alert('Username is not valid.');
            return
        }

        if (isRegistration && form.email) {
            if (!validateEmail(form.email)) {
                // CHANGE ALERT ON CASTOM
                alert('Email is not valid.');
                return
            }
        }

        let res: AuthResponse
        if (isRegistration) {
            res = await registerUser(form);
        } else {
            res = await loginUser(form);
        }

        if (res.error) {
            // CHANGE ALERT ON CASTOM
            alert(res.message);
        } else {
            if (isRegistration) {
                window.location.href = '/login';
            } else {
                window.location.href = '/app/wishlist';
            }
        }
    }

    return (
        <div className={styles.form}>
            {isRegistration && <h2>Registration</h2>}
            {!isRegistration && <h2>Login</h2>}

            <form className={styles.formInner}>
                <input
                    type="text"
                    name="username"
                    placeholder="Username"
                    required
                    value={form.username}
                    onChange={handleChange}
                />
                <input
                    type="password"
                    name="password"
                    placeholder="Password"
                    required
                    value={form.password}
                    onChange={handleChange}
                />
                {isRegistration && (
                    <input
                        type="email"
                        name="email"
                        placeholder="Email"
                        required
                        value={form.email}
                        onChange={handleChange}
                    />
                )}
                <Button handler={handleClick} text={isRegistration ? "Register" : "Login"} />
                {!isRegistration && (
                    <div className={styles.links}>
                        <Link to="/registration" className={styles.link}>Create account</Link>
                    </div>
                )}
            </form>
        </div>
    );
}
