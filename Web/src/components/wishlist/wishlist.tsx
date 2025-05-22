import { useState, useEffect } from 'react';
import { WishItem as WishItemType } from '../../api/wishlist';
import { getWishlist } from '../../api/wishlist';
import classes from './wishlist.module.css';
import wishItemImage from '../../assets/wish_not_found.jpg';
import { Link } from '@tanstack/react-router';


export function Wishlist() {
    const [wishlist, setWishlist] = useState<WishItemType[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        async function fetchWishlist() {
            try {
                const data = await getWishlist();
                setWishlist(data);
            } catch (err) {
                setError('Failed to load wishlist');
                console.error(err);
            } finally {
                setLoading(false);
            }
        }

        fetchWishlist();
    }, []);

    function WishItem(props: WishItemType) {
        const { id, image_url, name, price, created_at } = props;

        return (
            <Link
                to={`/app/$wishId`}
                params={{
                    wishId: id,
                }}
                className={classes.wishItem}>
                <div className={classes.wishItemImage}>
                    <img src={image_url ? image_url : wishItemImage} alt={name} />
                </div>
                <h3>{name}</h3>
                <div className={classes.wishItemInfo}>
                    <span className={classes.wishItemPrice}>{price ? price + ' $' : ''}</span>
                    <span className={classes.wishItemDate}>
                        {new Date(created_at).toDateString()}
                    </span>
                </div>
            </Link>
        );
    }

    if (loading) {
        return;
    }

    if (error) {
        return;
    }

    return (
        <section className={classes.wishlistSection}>
            <h2>Your Wishlist</h2>
            <div className={classes.wishlist}>
                {wishlist.map((item, index) => (
                    <WishItem key={index} {...item} />
                ))}
            </div>
        </section>
    );
}