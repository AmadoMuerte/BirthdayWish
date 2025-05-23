import { getTokenInfo } from "./token";

const API_URL = import.meta.env.VITE_API_URL;

export type WishItem = {
    id: number;
    image_url: string;
    image_name: string;
    link: string;
    name: string;
    price: number;
    created_at: Date;
    updated_at: Date;
}

export type WishItemReq = {
    price: number;
    link: string;
    image_data: string;
    image_type: string;
    name: string;
}

export async function getWishlist(token: string | null, path: string): Promise<WishItem[]> {
    try {
        const response = await fetch(API_URL + path, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + token
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

export async function addToWishlist(wishItem: WishItemReq) {
    const tokenInfo = getTokenInfo();

    try {
        const response = await fetch(API_URL + '/api/wishlist', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + tokenInfo.token
            },
            body: JSON.stringify(wishItem)
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        return data;

    } catch (error) {
        console.error('Error adding to wishlist:', error);
        throw error;
    }
}