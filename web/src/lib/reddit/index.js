/* https://www.reddit.com/r/redditdev/wiki/oauth2/quickstart/ */
export const state = '7A3F9E2C8B1D';
const TOKEN_URL = 'https://ssl.reddit.com/api/v1/access_token';

// connects and gets access token for username/password
export async function requestToken (code) {
    const response = await fetch('/api/v1/reddit/token', {
        method: 'POST',
        body: JSON.stringify({ code })
    }); 

    if (response.ok) {
        return [null, await response.json()];
    }

    return [new Error('Failed to login')];
}
