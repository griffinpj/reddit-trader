package auth

import (
	"context"
	"net/http"
)

// ContextKey for storing user claims in context
type ContextKey string

const UserClaimsKey ContextKey = "userClaims"

// RequireAuth is a middleware that validates JWT tokens
func (j *JWTManager) RequireNoAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from cookie
		token, err := j.GetTokenFromCookie(r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		// Validate token
		_, err = j.ValidateToken(token)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		
		http.Redirect(w, r, "/", http.StatusSeeOther);
	})
}


// RequireAuth is a middleware that validates JWT tokens
func (j *JWTManager) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from cookie
		token, err := j.GetTokenFromCookie(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther);
			return
		}

		// Validate token
		claims, err := j.ValidateToken(token)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther);
			return
		}

		// Add claims to context
		ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserClaims retrieves user claims from context
func GetUserClaims(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(UserClaimsKey).(*Claims)
	return claims, ok
}
