export interface JwtPayload {
    user_id: number;
    exp: number;
    iat: number;
}

export interface TokenInfo {
    hasToken: boolean;
    token: string | null;
    isValid: boolean;
    payload: JwtPayload | null;
}

const TOKEN_KEY = 'jwt_token';

export const setToken = (token: string): boolean => {
    if (!isTokenValid(token)) {
        console.error('Invalid token format');
        return false;
    }

    localStorage.setItem(TOKEN_KEY, token);
    return true;
};

export const getTokenInfo = (): TokenInfo => {
    const token = localStorage.getItem(TOKEN_KEY);

    if (!token) {
        return {
            hasToken: false,
            token: null,
            isValid: false,
            payload: null
        };
    }

    try {
        const payload = parseJwtPayload(token);
        const isValid = isTokenValid(token);

        return {
            hasToken: true,
            token,
            isValid,
            payload: isValid ? payload : null
        };
    } catch (e) {
        console.error(e);

        return {
            hasToken: false,
            token: null,
            isValid: false,
            payload: null
        };
    }
};

export const removeToken = (): void => {
    localStorage.removeItem(TOKEN_KEY);
};


const parseJwtPayload = (token: string): JwtPayload => {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(
        atob(base64)
            .split('')
            .map(c => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
            .join('')
    );

    return JSON.parse(jsonPayload);
};

const isTokenValid = (token: string | null): boolean => {
    if (!token) return false;

    try {
        const payload = parseJwtPayload(token);
        const now = Math.floor(Date.now() / 1000);
        return payload.exp > now && payload.iat <= now;
    } catch (e) {
        console.error(e);
        return false;
    }
};