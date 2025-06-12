import { NumberInput, TextInput } from '@mantine/core'
import {
  CreateWishData,
  UpdateWishData,
  Wish
} from '@/entities/wish/model/types'
import { useForm } from '@mantine/form'
import { useCreateWishApi, useUpdateWishApi } from '@/entities/wish/model/hooks'

import styles from './wish-form.module.css'

interface WishFormProps {
  initialValues?: Wish
}

export const WishForm: React.FC<WishFormProps> = ({ initialValues }) => {
  const form = useForm({
    mode: 'uncontrolled',
    initialValues
  })
  const { mutate: update } = useUpdateWishApi({})
  const { mutate: create } = useCreateWishApi({})

  const handleSubmit = (values: CreateWishData | UpdateWishData) => {
    if (initialValues) {
      update({ body: values as UpdateWishData, id: initialValues.id })
    } else {
      create({ body: values as CreateWishData })
    }
  }

  return (
    <form
      className={styles['wish-form']}
      onSubmit={form.onSubmit((values) => handleSubmit(values))}
    >
      <div>
        <TextInput
          withAsterisk
          label='Name'
          key={form.key('name')}
          {...form.getInputProps('name')}
        />
        <NumberInput
          withAsterisk
          label='Price'
          key={form.key('price')}
          {...form.getInputProps('price')}
        />
        <TextInput
          withAsterisk
          label='Link'
          key={form.key('link')}
          {...form.getInputProps('link')}
        />
      </div>
      <div>

      </div>
    </form>
  )
}
