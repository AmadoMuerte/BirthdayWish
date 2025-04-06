const API_URL = 'http://localhost:3030';

export type Form = {
    username: string;
    password: string;
    email?: string;
};

export async function loginUser(form: Form) {
    const response = await fetch(`${API_URL}/auth/login`, {
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
