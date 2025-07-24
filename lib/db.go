package lib

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/jackc/pgx/v5/pgxpool"

	Config "rtrade/config"
)

func Database () * pgxpool.Pool {
	var config *Config.Config;
	var connectionStr string;
	var pool * pgxpool.Pool;
	var err error;

	config, err = Config.Load();
	if (err != nil) {
		log.Fatal("Failed to load config for db");
	}

	connectionStr = "postgres://" + config.Store.User + ":" + config.Store.Password + "@" + config.Store.Host + "/" + config.Store.Name + "?sslmode=" + config.Store.Ssl;
	pool, err = pgxpool.New(context.Background(), connectionStr);
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err);
	}

	// defer conn.Close(context.Background());
	
	return pool
}

// SeedAdminUser creates an initial admin user in the database
func seedAdminUser(db *sql.DB) error {
    // Admin user details
    email := "admin@example.com"
    username := "admin"
    password := "Admin123!" // Change this to a secure password
    firstName := "Admin"
    lastName := "User"
    displayName := "Administrator"

    // Check if admin user already exists
    var exists bool
    err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 OR username = $2)", email, username).Scan(&exists)
    if err != nil {
        return fmt.Errorf("error checking existing admin user: %w", err)
    }

    if exists {
        fmt.Println("Admin user already exists, skipping seed")
        return nil
    }

    // Generate salt
    salt := make([]byte, 16)
    if _, err := rand.Read(salt); err != nil {
        return fmt.Errorf("error generating salt: %w", err)
    }
    saltHex := hex.EncodeToString(salt)

    // Hash password with salt
    // Note: bcrypt internally handles salting, so we're combining our salt with password for extra security
    combinedPassword := []byte(password + saltHex)

    hashedPassword, err := bcrypt.GenerateFromPassword(combinedPassword, bcrypt.DefaultCost)
    if err != nil {
        return fmt.Errorf("error hashing password: %w", err)
    }

    // Insert admin user
    query := `
        INSERT INTO users (
            email, username, password_hash, password_salt,
            first_name, last_name, display_name,
            is_active, is_verified, role,
            email_verified_at, password_changed_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
        )`

    _, err = db.Exec(
        query,
        email,
        username,
        string(hashedPassword),
        saltHex,
        firstName,
        lastName,
        displayName,
        true,  // is_active
        true,  // is_verified
        "admin", // role
        time.Now(), // email_verified_at
        time.Now(), // password_changed_at
    )

    if err != nil {
        return fmt.Errorf("error inserting admin user: %w", err)
    }

    fmt.Println("Admin user created successfully")
    fmt.Printf("Email: %s\n", email)
    fmt.Printf("Username: %s\n", username)
    fmt.Printf("Password: %s (Please change this immediately!)\n", password)

    return nil
}

func createTables (db *sql.DB) sql.Result {
	result, err := db.Exec(`
		CREATE TABLE users (
    -- Primary key
    id BIGSERIAL PRIMARY KEY,

    -- Authentication fields
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    password_salt VARCHAR(255) NOT NULL,

    -- Profile information
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    display_name VARCHAR(100),
    bio TEXT,
    avatar_url VARCHAR(500),

    -- Contact information
    phone_number VARCHAR(20),

    -- Account status
    is_active BOOLEAN DEFAULT true,
    is_verified BOOLEAN DEFAULT false,
    role VARCHAR(50) DEFAULT 'user',

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_login_at TIMESTAMP WITH TIME ZONE,

    -- Additional security
    email_verified_at TIMESTAMP WITH TIME ZONE,
    password_changed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CONSTRAINT email_format CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    CONSTRAINT username_format CHECK (username ~* '^[a-zA-Z0-9_]{3,50}$')
);

-- Indexes for performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_is_active ON users(is_active) WHERE is_active = true;

-- Trigger to automatically update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
	`);
	if err != nil {
		log.Fatal("Failed to create tables")
	}

	return result;
}
