package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Config holds JWT configuration
type Config struct {
	SecretKey      string
	TokenExpiry    time.Duration
	RefreshExpiry  time.Duration
	CookieName     string
	CookieDomain   string
	CookieSecure   bool
	CookieHTTPOnly bool
	CookieSameSite http.SameSite
}

// Claims represents the JWT claims
type Claims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JWTManager handles JWT operations
type JWTManager struct {
	config Config
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(config Config) *JWTManager {
	// Set defaults if not provided
	if config.CookieName == "" {
		config.CookieName = "auth_token"
	}
	if config.TokenExpiry == 0 {
		config.TokenExpiry = 24 * time.Hour
	}
	if config.CookieSameSite == 0 {
		config.CookieSameSite = http.SameSiteLaxMode
	}
	config.CookieHTTPOnly = true // Always true for security

	return &JWTManager{
		config: config,
	}
}

// GenerateToken creates a new JWT token
func (j *JWTManager) GenerateToken(userID, email, username string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Email:    email,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.config.TokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.config.SecretKey))
}

// ValidateToken validates and parses a JWT token
func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.config.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// SetAuthCookie sets the JWT token as an HTTP cookie
func (j *JWTManager) SetAuthCookie(w http.ResponseWriter, token string) {
	cookie := &http.Cookie{
		Name:     j.config.CookieName,
		Value:    token,
		Path:     "/",
		Domain:   j.config.CookieDomain,
		Expires:  time.Now().Add(j.config.TokenExpiry),
		MaxAge:   int(j.config.TokenExpiry.Seconds()),
		Secure:   j.config.CookieSecure,
		HttpOnly: j.config.CookieHTTPOnly,
		SameSite: j.config.CookieSameSite,
	}
	http.SetCookie(w, cookie)
}

// ClearAuthCookie removes the auth cookie
func (j *JWTManager) ClearAuthCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     j.config.CookieName,
		Value:    "",
		Path:     "/",
		Domain:   j.config.CookieDomain,
		Expires:  time.Now().Add(-1 * time.Hour),
		MaxAge:   -1,
		Secure:   j.config.CookieSecure,
		HttpOnly: j.config.CookieHTTPOnly,
		SameSite: j.config.CookieSameSite,
	}
	http.SetCookie(w, cookie)
}

// GetTokenFromCookie extracts the token from the request cookie
func (j *JWTManager) GetTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(j.config.CookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

