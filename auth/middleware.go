package auth

import (
	"context"
	"log"
	"net/http"
	"rtrade/db"

	"github.com/jackc/pgx/v5"
)

// ContextKey for storing user claims in context
type ContextKey string

const UserClaimsKey ContextKey = "userClaims"

// TODO will the jwt just expire? 
// Are we supposed to renew on new requests?

// TODO cleanup attaching to context

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
			http.Error(w, err.Error(), http.StatusUnauthorized);
			return
		}

		// Validate token
		claims, err := j.ValidateToken(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized);
			return
		}

		// Add claims to context
		ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
		
		conn, err := j.config.Pool.Acquire(context.Background())
		if err != nil {
			log.Println(err);
			http.Error(w, err.Error(), http.StatusInternalServerError);
			return;
		}

		// Return the aquired db connection back to the connection pool from which it came!
		defer conn.Release();
		q := db.New(conn);
		
		log.Println(claims.UserID)
		if err != nil {
			log.Println(err);
			http.Error(w, err.Error(), http.StatusInternalServerError);
			return;
		}

		// Get User and attach to context
		var user db.User
		user, err = q.GetUserById(context.Background(), claims.UserID);
		if err != nil && err != pgx.ErrNoRows {
			log.Println(err);
			http.Error(w, err.Error(), http.StatusInternalServerError);
			return;
		}

		if err == pgx.ErrNoRows {
			log.Println(err);
			http.Error(w, err.Error(), http.StatusInternalServerError);
			return;
		}

		ctx = context.WithValue(ctx, "user", user);
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAuth is a middleware that validates JWT tokens
func (j *JWTManager) RequireAuthRedir(next http.Handler) http.Handler {
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
		ctx := context.WithValue(r.Context(), UserClaimsKey, claims);
		conn, err := j.config.Pool.Acquire(context.Background());
		if err != nil {
			log.Println(err);
			http.Error(w, err.Error(), http.StatusInternalServerError);
			return;
		}

		// Return the aquired db connection back to the connection pool from which it came!
		defer conn.Release();
		q := db.New(conn);
		
		log.Println(claims.UserID)
		if err != nil {
			log.Println(err);
			http.Error(w, err.Error(), http.StatusInternalServerError);
			return;
		}

		// Get User and attach to context
		var user db.User
		user, err = q.GetUserById(context.Background(), claims.UserID);
		if err != nil && err != pgx.ErrNoRows {
			log.Println(err);
			http.Error(w, err.Error(), http.StatusInternalServerError);
			return;
		}

		if err == pgx.ErrNoRows {
			log.Println(err);
			http.Error(w, err.Error(), http.StatusInternalServerError);
			return;
		}

		ctx = context.WithValue(ctx, "user", user);

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserClaims retrieves user claims from context
func GetUserClaims(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(UserClaimsKey).(*Claims)
	return claims, ok
}
