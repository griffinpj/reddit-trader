# reddit-trader

## Usage 🔥
```
    go get
    go mod tidy
    go run .
```

### Requirements 🛠️

- `postgresql`
- `sqlc`
> Update schema and queries in `config/db/...`.
> Use, `sqlc generate` to update queries and generate type-safe code.

### Environment

Ensure your environment variables are set

```
APP_PORT=3333

# Reddit
REDDIT_CLIENT_ID=
REDDIT_CLIENT_SECRET=
REDDIT_TOKEN_URL=https://ssl.reddit.com/api/v1/access_token
REDDIT_API_URL=https://oauth.reddit.com/api/v1

# db
DB_USER=
DB_PASSWORD=
DB_HOST=
DB_NAME=
DB_SSL=
```

