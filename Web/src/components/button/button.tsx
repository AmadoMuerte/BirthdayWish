import styles from './button.module.css'


type ButtonProps = {
    text: string,
    handler: () => void
}

export default function Button(props: ButtonProps) {

    return (
        <button onClick={props.handler} type="submit" className={styles.btn}>{props.text}</button>
    )
}