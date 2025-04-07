import { getTokenInfo, removeToken, setToken, TokenInfo } from "./token";

const API_URL = import.meta.env.VITE_API_URL;

export interface AuthResponse {
    error: boolean;
    message: string;
    tokenInfo?: TokenInfo
}

export type Form = {
    username: string;
    password: string;
    email?: string;
};

export async function loginUser(form: Form): Promise<AuthResponse> {
    try {
        const response = await fetch(`${API_URL}/auth/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(form),
        });

        if (!response.ok) {
            removeToken();

            let errorMessage = response.statusText;
            try {
                const errorData = await response.json();
                errorMessage = errorData.message || errorMessage;
            } catch (e) {
                console.error('Failed to parse error response', e);
            }

            return {
                error: true,
                message: errorMessage
            };
        }

        const data = await response.json();

        if (!data.token) {
            return {
                error: true,
                message: 'Token not received from server'
            };
        }

        const tokenSaved = setToken(data.token);
        if (!tokenSaved) {
            return {
                error: true,
                message: 'Failed to save authentication token'
            };
        }
        const tokenInfo = getTokenInfo();

        return {
            error: false,
            message: 'Login successful',
            tokenInfo
        };

    } catch (error) {
        console.error('Login error:', error);
        return {
            error: true,
            message: error instanceof Error ? error.message : 'An unknown error occurred'
        };
    }
}

export async function registerUser(form: Form) {
    const response = await fetch(`${API_URL}/auth/sign_up`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(form),
    });

    if (response.ok) {
        return await response.json();
    } else {
        throw new Error(response.statusText);
    }
}
