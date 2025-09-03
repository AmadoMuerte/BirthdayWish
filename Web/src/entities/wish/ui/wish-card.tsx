import { Link } from '@tanstack/react-router'

import { Wish } from '../model/types'
import wishItemImage from '@/assets/wish_not_found.jpg'

import styles from './wish-card.module.css'

export interface WishCardProps {
  wish: Wish
}

export const WishCard: React.FC<WishCardProps> = ({ wish }) => {
  return (
    <Link
      to={`/app/wishes/$id`}
      params={{ id: wish.id }}
      className={styles['wish-card']}
    >
      <div className={styles['wish-card__image']}>
        <img
          src={wish.image ? wish.image : wishItemImage}
          alt={wish.name}
        />
      </div>
      <h3 className={styles['wish-card__name']}>{wish.name}</h3>
      <div className={styles['wish-card__info']}>
        <span className={styles['wish-card__price']}>
          {wish.price ? wish.price + ' $' : ''}
        </span>
        <span className={styles['wish-card__date']}>
          {new Date(wish.created_at).toDateString()}
        </span>
      </div>
    </Link>
  )
}
