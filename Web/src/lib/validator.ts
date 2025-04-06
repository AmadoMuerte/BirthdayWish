export default function validateEmail(email: string): boolean {
    const re = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
    return re.test(email.toLowerCase());
}

export function validatePassword(password: string): boolean {
    const re = /^[a-zA-Z0-9!@#$%^&*()\-_=+]{8,20}$/;
    return re.test(password);
}

export function validateUsername(username: string): boolean {
    const re = /^[a-zA-Z0-9_-]{3,20}$/;
    return re.test(username);
}
