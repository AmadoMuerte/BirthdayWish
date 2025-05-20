import classes from './wishlist.module.css'

type Wishitem = {
    image: string
    link: string
    name: string
    price: number
    createdAt: Date
}

const mocks = [
    {
        image: 'https://cdn.theatlantic.com/thumbor/zsL6Y0Fivx1gSWaLKiHS6A5F7LQ=/0x0:4800x2700/976x549/media/img/mt/2022/12/What_Gifts_Say/original.jpg',
        link: 'https://www.google.com',
        name: 'Name',
        price: 100,
        createdAt: new Date()
    },
    {
        image: 'https://cdn.theatlantic.com/thumbor/zsL6Y0Fivx1gSWaLKiHS6A5F7LQ=/0x0:4800x2700/976x549/media/img/mt/2022/12/What_Gifts_Say/original.jpg',
        link: 'https://www.google.com',
        name: 'Name',
        price: 100,
        createdAt: new Date()
    },
    {
        image: 'https://cdn.theatlantic.com/thumbor/zsL6Y0Fivx1gSWaLKiHS6A5F7LQ=/0x0:4800x2700/976x549/media/img/mt/2022/12/What_Gifts_Say/original.jpg',
        link: 'https://www.google.com',
        name: 'Name',
        price: 100,
        createdAt: new Date()
    },
    {
        image: 'https://cdn.theatlantic.com/thumbor/zsL6Y0Fivx1gSWaLKiHS6A5F7LQ=/0x0:4800x2700/976x549/media/img/mt/2022/12/What_Gifts_Say/original.jpg',
        link: 'https://www.google.com',
        name: 'Name',
        price: 100,
        createdAt: new Date()
    },
    {
        image: 'https://cdn.theatlantic.com/thumbor/zsL6Y0Fivx1gSWaLKiHS6A5F7LQ=/0x0:4800x2700/976x549/media/img/mt/2022/12/What_Gifts_Say/original.jpg',
        link: 'https://www.google.com',
        name: 'Name',
        price: 100,
        createdAt: new Date()
    },
    {
        image: 'https://cdn.theatlantic.com/thumbor/zsL6Y0Fivx1gSWaLKiHS6A5F7LQ=/0x0:4800x2700/976x549/media/img/mt/2022/12/What_Gifts_Say/original.jpg',
        link: 'https://www.google.com',
        name: 'Name',
        price: 100,
        createdAt: new Date()
    },
]

export function Wishlist() {

    function WishItem(props: Wishitem) {
        const { image, link, name, price, createdAt } = props

        return (
            <div className={classes.wishItem}>
                <div className={classes.wishItemImage}>
                    <img src={image} alt="" />
                </div>
                <h3>{name}</h3>
                <div className={classes.wishItemInfo}>
                    <span className={classes.wishItemPrice}>${price}</span>
                    <span className={classes.wishItemDate}>{createdAt.toDateString()}</span>
                </div>
            </div>
        )
    }


    return (
        <section className={classes.wishlistSection}>
            <h2>Your Wishlist</h2>
            <div className={classes.wishlist}>
                {mocks.map((wishitem, index) => (
                    <WishItem key={index} {...wishitem} />
                ))}
            </div>
        </section>
    )
}