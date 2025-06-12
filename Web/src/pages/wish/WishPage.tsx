import { Button, Container } from '@mantine/core'
import { Link, useRouter } from '@tanstack/react-router'
import { useGetWishApi } from '@/entities/wish/model/hooks'

import styles from './WishPage.module.css'

interface WishProps {
  id: number
}

export const WishPage: React.FC<WishProps> = ({ id }) => {
  const router = useRouter()
  const { data, isLoading, error } = useGetWishApi({ id })

  function handleBack() {
    router.history.back()
  }

  if (isLoading) {
    return <div>Loading...</div>
  }

  if (error) {
    return <div>Error: {error.message}</div>
  }

  return (
    <Container>
      <div className={styles['wish-page']}>
        <div className={styles['wish-page__header']}>
          <Button onClick={handleBack}>Вернуться</Button>
          <Link
            className={styles['wish-info__link']}
            to={`/app/wishes/${data?.id}/edit`}
          >
            <Button>Перейти к товару</Button>
          </Link>
        </div>
        <img
          className={styles['wish-page__img']}
          src={data?.image_url}
          alt={data?.name}
        />
        <div className={styles['wish-page__info']}>
          <div className={styles['wish-info']}>
            <h3 className={styles['wish-info__name']}>{data?.name}</h3>

            <div className={styles['wish-info__price']}>
              <div className={styles['wish-price']}>
                <p className={styles['wish-price__label']}>Цена:</p>
                <p className={styles['wish-price__value']}>
                  {data?.price} $
                </p>
              </div>
            </div>

            <Link className={styles['wish-info__link']} to={data?.link}>
              <Button>Перейти к товару</Button>
            </Link>
          </div>
        </div>
      </div>
    </Container>
  )
}
