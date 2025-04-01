import { Link } from "@tanstack/react-router";
import routes from "../../routes/routes.json";
import styles from './nav.module.css'

const iconModules = import.meta.glob('../../assets/icons/*.svg', { eager: true });

type Route = {
    name: string;
    path: string;
    icon?: string
}

function renderRoutes() {
    return routes.map((route: Route) => {
        const iconPath = route.icon ? `../../assets/icons/${route.icon}` : null;
        const iconSrc = iconPath ? (iconModules[iconPath] as { default: string })?.default : null;

        return (
            <li key={route.name}>
                <div className={styles.navItem}>
                    {iconSrc && <img src={iconSrc} alt={route.name} />}
                    <Link to={route.path}>
                        {route.name}
                    </Link>
                </div>
            </li>
        );
    })
}

export default function Nav() {
    return (
        <nav className={styles.nav}>
            <ul>
                {renderRoutes()}
            </ul>
        </nav>
    );
}