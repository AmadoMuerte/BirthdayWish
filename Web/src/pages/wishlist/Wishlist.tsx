import { WishCardList } from '@/entities/wish/ui/wish-card-list'
import { useGetAllWishesApi } from '@/entities/wish/model/hooks'

import styles from './Wishlist.module.css'

export function Wishlist() {
  const { data: wishes, error, isLoading } = useGetAllWishesApi({})

  if (isLoading) {
    return <>loading...</>
  }

  if (error) {
    return <>Error :c</>
  }

  if (!wishes || wishes.length === 0) {
    return <>No wishes found</>
  }

  return (
    <section className={styles.wishlistSection}>
      <h2>Aboba's Wishlist</h2>
      <WishCardList wishes={wishes} />
    </section>
  )
}
