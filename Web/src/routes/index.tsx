import { createFileRoute } from '@tanstack/react-router'
import Hero from '../components/landing/hero/hero'
import Features from '../components/landing/features/features'

export const Route = createFileRoute('/')({
    component: Index,
})

function Index() {
    return (
        <div>
            <Hero />
            <Features />
        </div>
    )
}