import styles from './button.module.css'


type ButtonProps = {
    text: string,
    handler: (e: React.MouseEvent<HTMLButtonElement>) => void
}

export default function Button(props: ButtonProps) {

    return (
        <button onClick={(e) => props.handler(e)} type="submit" className={styles.btn}>{props.text}</button>
    )
}