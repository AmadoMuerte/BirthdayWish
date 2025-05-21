import { getTokenInfo } from "./token";

const API_URL = import.meta.env.VITE_API_URL;

export type WishItem = {
    image_url: string;
    image_name: string;
    link: string;
    name: string;
    price: number;
    created_at: Date;
    updated_at: Date;
}

export async function getWishlist(): Promise<WishItem[]> {
    const tokenInfo = getTokenInfo();

    try {
        const response = await fetch(API_URL + '/api/wishlist/' + tokenInfo.payload?.user_id, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + tokenInfo.token
            }
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        return data as WishItem[];
    } catch (error) {
        console.error('Error fetching wishlist:', error);
        return [];
    }
}