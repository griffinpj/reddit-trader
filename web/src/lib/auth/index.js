
export async function logout () {
    const response = await fetch('/auth/logout', {
        method: 'POST'
    }); 

    if (response.ok) {
        window.location.reload();
    }

    return new Error('Failed to logout');
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
