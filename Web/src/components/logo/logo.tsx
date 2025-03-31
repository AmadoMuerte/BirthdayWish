import { Link } from "@tanstack/react-router"
import styles from "./logo.module.css"

export default function Logo() {

    return (
        <Link to="/" className={styles.logo}>
            BirthdayWish
        </Link>
    )
}