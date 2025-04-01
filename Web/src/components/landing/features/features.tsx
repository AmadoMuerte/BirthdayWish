import styles from './features.module.css'

export default function Features() {

    const features = [
        {
            name: 'feature 1',
            description: 'description 1',
            image: 'icon 1'
        },
        {
            name: 'feature 2',
            description: 'description 1',
            image: 'icon 1'
        },
        {
            name: 'feature 3',
            description: 'description 1',
            image: 'icon 1'
        }
    ]

    return (
        <section id="features" className={styles.features}>
            <div className="container">
                <div className={styles.featresInner}>
                    <h2 className='sectionTitle'>Features</h2>
                    <div className={styles.featuresList}>
                        {features.map((feature, index) => (
                            <div key={index} className={styles.feature}>
                                <h3>{feature.name}</h3>
                                <p>{feature.description}</p>
                            </div>
                        ))}
                    </div>
                </div>
            </div>
        </section>
    )
} 