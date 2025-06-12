import { Wish } from '../model/types'
import { WishCard } from './wish-card'

import styles from './wish-card-list.module.css'

export interface WishCardListProps {
  wishes: Wish[]
}

export const WishCardList: React.FC<WishCardListProps> = ({ wishes }) => {
  return (
    <div className={styles['wish-card-list']}>
      {wishes.map((wish) => (
        <WishCard key={wish.id} wish={wish} />
      ))}
    </div>
  )
}
