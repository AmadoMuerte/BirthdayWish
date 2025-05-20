import { Navbar } from "../nav_bar/nav_bar";

type AppWrapperProps = {
    content: React.ReactNode
}

export function AppWrapper(props: AppWrapperProps) {

    return (
        <div>
            <Navbar />
            <div className="container">
                {props.content}
            </div>
        </div>
    )
}