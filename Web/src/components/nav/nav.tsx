import { Link } from "@tanstack/react-router";
import routes from "../../routes/routes.json";
import styles from './nav.module.css'

type Route = {
    name: string;
    path: string;
}

export default function Nav() {

    return (
        <nav className={styles.nav}>
            <ul>
                {routes.map((route: Route) => (
                    <li key={route.name}>
                        <Link to={route.path}>
                            {route.name}
                        </Link>{' '}
                    </li>
                ))}
            </ul>
        </nav>
    )
}