import { createFileRoute, useSearch } from '@tanstack/react-router'
import { Wishlist } from '../pages/wishlist/Wishlist'
import { getWishlist, WishItem as WishItemType } from '../shared/api/wishlist';
import { useEffect, useState } from 'react';

export const Route = createFileRoute('/sharelist')({
  component: RouteComponent,
})

function RouteComponent() {
  const search: { tkn: string } = useSearch({ strict: false })

  const [wishlist, setWishlist] = useState<WishItemType[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);


  useEffect(() => {
    async function fetchWishlist() {
      try {
        const data = await getWishlist(search.tkn, '/auth/get_wishlist');
        setWishlist(data);
        if (data.length === 0) {
          setError('Wishlist is empty');
        }
      } catch (err) {
        setError('Failed to load wishlist');
        console.error(err);
      } finally {
        setLoading(false);
      }
    }

    fetchWishlist();
  }, []);

  if (loading) {
    return
  }

  if (error) {
    window.location.href = '/';
  }

  return (
    <div className="container">
      <Wishlist wishlist={wishlist} loading={loading} error={error} />
    </div>
  )
}
