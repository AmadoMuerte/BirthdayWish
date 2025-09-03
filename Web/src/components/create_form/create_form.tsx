import { useState, useCallback, ChangeEvent, DragEvent, FormEvent } from 'react'
import { useRouter } from '@tanstack/react-router'
import { addToWishlist } from '../../shared/api/wishlist'
import styles from './create_form.module.css'

interface FormData {
  name: string
  price: number | undefined
  link: string
  image: string | null
  image_type: string | null
}

interface CreateFormProps {
  initialData?: FormData
}

export const CreateForm: React.FC<CreateFormProps> = ({ initialData }) => {
  const router = useRouter()

  const [imagePreview, setImagePreview] = useState<string>('')
  const [formData, setFormData] = useState<FormData>(
    initialData || {
      name: '',
      price: undefined,
      link: '',
      image: null,
      image_type: null
    }
  )

  const handleInputChange = (e: ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target
    setFormData((prev) => ({
      ...prev,
      [name]: value
    }))
  }

  const handleImageChange = useCallback((e: ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      const reader = new FileReader()
      reader.onloadend = () => {
        const result = reader.result as string
        setFormData((prev) => ({
          ...prev,
          image: result.split(',')[1],
          image_type: file.type,
        }))
        setImagePreview(result)
      }
      reader.readAsDataURL(file)
    }
  }, [])

  const handleDragOver = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault()
  }

  const handleDrop = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault()
    const file = e.dataTransfer.files?.[0]
    if (file) {
      const reader = new FileReader()
      reader.onloadend = () => {
        const result = reader.result as string
        setFormData((prev) => ({
          ...prev,
          image: result.split(',')[1],
          image_type: file.type
        }))
        setImagePreview(result)
      }
      reader.readAsDataURL(file)
    }
  }

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault()

    if (!formData.name.trim()) {
      alert('Please enter a valid name')
      return
    }

    if (formData.price && isNaN(Number(formData.price))) {
      alert('Please enter a valid price')
      return
    }

    const res = addToWishlist({
      price: Number(formData.price),
      link: formData.link,
      image: formData.image as string,
      image_type: formData.image_type as string,
      name: formData.name
    })

    if (!res) {
      alert('Failed to add to wishlist')
      return
    }

    router.navigate({
      to: '/app/wishlist'
    })
  }

  return (
    <section className={styles.sectionForm}>
      <div className={styles.container}>
        <h1 className={styles.title}>Add a product</h1>
        <p className={styles.description}>
          Fill out the form to add the desired gift, only the name field is
          required <br />
          <span>you fill out the remaining fields as desired</span>
        </p>
        <form onSubmit={handleSubmit}>
          <div className={styles.form_wrapper}>
            <div className={styles.form_inputWrapper}>
              <div className={styles.formGroup}>
                <label className={styles.label}>Name</label>
                <input
                  type='text'
                  name='name'
                  value={formData.name}
                  onChange={handleInputChange}
                  className={styles.input}
                  placeholder='Product name'
                />
              </div>
              <div className={styles.formGroup}>
                <label className={styles.label}>Price</label>
                <input
                  type='text'
                  name='price'
                  value={formData.price}
                  onChange={handleInputChange}
                  className={styles.input}
                  placeholder='Product price'
                />
              </div>
              <div className={styles.formGroup}>
                <label className={styles.label}>Product link</label>
                <input
                  type='url'
                  name='link'
                  value={formData.link}
                  onChange={handleInputChange}
                  className={styles.input}
                  placeholder='Product link'
                />
              </div>
            </div>
            <div className={styles.uploadSection}>
              <div
                className={styles.uploadArea}
                onDragOver={handleDragOver}
                onDrop={handleDrop}
                onClick={() => document.getElementById('file-upload')?.click()}
              >
                {imagePreview ? (
                  <img
                    src={imagePreview}
                    alt='Preview'
                    className={styles.imagePreview}
                  />
                ) : (
                  <p className={styles.uploadText}>
                    <span>Upload a picture</span>
                    <br /> Drag and drop or click to upload
                  </p>
                )}
                <input
                  id='file-upload'
                  type='file'
                  accept='image/*'
                  onChange={handleImageChange}
                  className={styles.fileInput}
                />
              </div>
            </div>
          </div>
          <button type='submit' className={styles.submitButton}>
            Add to Wishlist
          </button>
        </form>
      </div>
    </section>
  )
}
