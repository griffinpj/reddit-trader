// https://www.reddit.com/r/redditdev/wiki/oauth2/quickstart/

package config

import (
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Config struct {
	Application AppConfig
	Reddit RedditConfig
	Store DBConfig
}

type AppConfig struct {
	Port string
}

type RedditConfig struct {
	ClientId string
	ClientSecret string
	TokenURL string
	ApiURL string
}

type DBConfig struct {
	User string
	Password string
	Host string
	Name string
	Ssl string
}

type Env struct {
	Pool * pgxpool.Pool
}

// Return a pointer to the original config to avoid making copies
func Load () (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables");
	}
	
	// Creating a config pointer directly
	cfg := &Config {
		Application: AppConfig {
			Port: getEnv("APP_PORT", "3333"), // need commas here for structs ...
		},
		Store: DBConfig {
			User: getEnv("DB_USER", ""),
			Password: getEnv("DB_PASSWORD", ""),
			Host: getEnv("DB_HOST", ""),
			Name: getEnv("DB_NAME", ""),
			Ssl: getEnv("DB_SSL", "verify-full"),
		},
		Reddit: RedditConfig {
			ClientId: getEnv("REDDIT_CLIENT_ID", ""),
            ClientSecret: getEnv("REDDIT_CLIENT_SECRET", ""),
            TokenURL: getEnv("REDDIT_TOKEN_URL", "https://ssl.reddit.com/api/v1/access_token"),
            ApiURL: getEnv("REDDIT_API_URL", "https://oauth.reddit.com/api/v1"),
		},
	}

	if cfg.Reddit.ClientId == "" || cfg.Reddit.ClientSecret == "" {
		return nil, log.New(os.Stderr, "", 0).Output(0, "REDDIT_CLIENT_ID and REDDIT_CLIENT_SECRET are required");
	}

	return cfg, nil
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value;
    }

    return fallback;
}
