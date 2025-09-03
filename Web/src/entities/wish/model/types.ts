export interface Wish {
  id: number;
  user_id: number;
  link: string;
  price: number;
  image: string;
  image_name: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export type CreateWishData = Omit<Wish, 'id' | 'user_id' | 'created_at' | 'updated_at'>;
export type UpdateWishData = Partial<CreateWishData> & { id: number };