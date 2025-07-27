
export async function logout () {
    const response = await fetch('/auth/logout', {
        method: 'POST'
    }); 

    if (response.ok) {
        return;
    }

    return new Error('Failed to logout');
}

export async function register (values) {
    const response = await fetch('/auth/register', {
        method: 'POST',
        body: JSON.stringify(values)
    }); 

    if (response.ok) {
        const data = await response.json();
        return [null, data]; 
    }

    return [new Error('Failed to login')];
}

export async function login (username, password) {
    const response = await fetch('/auth/login', {
        method: 'POST',
        body: JSON.stringify({ username, password })
    }); 

    if (response.ok) {
        return response.json();
    }

    return new Error('Failed to login');
}
