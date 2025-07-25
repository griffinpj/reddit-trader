package routers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"rtrade/config"
	"rtrade/db"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func Auth (env *config.Env) chi.Router {
	r := chi.NewRouter();

	type ErrorResponse struct {
		Message string `json:"message"`
	}
	
	r.Post("/logout", func (w http.ResponseWriter, r *http.Request) {
		env.Jwt.ClearAuthCookie(w);
		w.Header().Set("Content-Type", "application/json");
		
		json.NewEncoder(w).Encode(ErrorResponse {
			Message: "Logout successful",
		});
	});
	
	r.Post("/login", func (w http.ResponseWriter, r *http.Request) {
		var LoginData struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		decoder := json.NewDecoder(r.Body);
		err := decoder.Decode(&LoginData)

		conn, err := env.Pool.Acquire(context.Background())
		if err != nil {
			log.Println(err);
			http.Error(w, err.Error(), http.StatusInternalServerError);
			return;
		}

		// Return the aquired db connection back to the connection pool from which it came!
		defer conn.Release();

		q := db.New(conn);
	
		var user db.User
		user, err = q.GetUser(context.Background(), LoginData.Username);
		if err != nil && err != pgx.ErrNoRows {
			log.Println(err);
			http.Error(w, err.Error(), http.StatusInternalServerError);
			return;
		}

		if err == pgx.ErrNoRows {
			w.Header().Set("Content-Type", "application/json");
			
			json.NewEncoder(w).Encode(ErrorResponse {
				Message: "Username does not exist",
			});
			return;
		}

		combinedPassword := []byte(LoginData.Password + user.PasswordSalt)
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), combinedPassword)
		if err != nil {
			log.Println(err);
			w.Header().Set("Content-Type", "application/json");
			json.NewEncoder(w).Encode(ErrorResponse {
				Message: "Username or Password is incorrect",
			});
			return;
		}

		// TODO verified password
		// TODO set auth cookie
		token, err := env.Jwt.GenerateToken(string(user.ID), user.Email, user.Username);
		if err != nil {
			log.Println(err);
			http.Error(w, err.Error(), http.StatusInternalServerError);
			return;
		}

		env.Jwt.SetAuthCookie(w, token);
		w.Header().Set("Content-Type", "application/json");
		json.NewEncoder(w).Encode(user);
	});

	return r;
}

