package routers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"rtrade/config"
	"rtrade/db"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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

	r.Post("/register", func (w http.ResponseWriter, r *http.Request) {
		var RegisterData struct {
			FirstName string `json:"first_name"`
			LastName string `json:"last_name"`
			Email string `json:"email"`
			Username string `json:"username"`
			Password string `json:"password"`
			Confirm string `json:"password-confirm"`
		}

		decoder := json.NewDecoder(r.Body);
		err := decoder.Decode(&RegisterData)

		if RegisterData.Password != RegisterData.Confirm {
			w.Header().Set("Content-Type", "application/json");
			
			json.NewEncoder(w).Encode(ErrorResponse {
				Message: "Passwords do not match",
			});
			return;
		}

		// Generate salt
		salt := make([]byte, 16)
		if _, err := rand.Read(salt); err != nil {
			log.Println(err);
			http.Error(w, err.Error(), http.StatusInternalServerError);
			return;
		}

		saltHex := hex.EncodeToString(salt)

		// Hash password with salt
		// Note: bcrypt internally handles salting, so we're combining our salt with password for extra security
		combinedPassword := []byte(RegisterData.Password + saltHex)

		hashedPassword, err := bcrypt.GenerateFromPassword(combinedPassword, bcrypt.DefaultCost)
		if err != nil {
			log.Println(err);
			http.Error(w, err.Error(), http.StatusInternalServerError);
			return;
		}

		conn, err := env.Pool.Acquire(context.Background())
		if err != nil {
			log.Println(err);
			http.Error(w, err.Error(), http.StatusInternalServerError);
			return;
		}

		// Return the aquired db connection back to the connection pool from which it came!
		defer conn.Release();
		q := db.New(conn);

		alreadyExists, err := q.UserExists(context.Background(), db.UserExistsParams{
			Username: RegisterData.Username,
			Email: RegisterData.Email,
		})
		if err != nil {
			log.Println(err);
			http.Error(w, err.Error(), http.StatusInternalServerError);
			return;
		}

		if alreadyExists {
			w.Header().Set("Content-Type", "application/json");
			
			json.NewEncoder(w).Encode(ErrorResponse {
				Message: "Username or Email is already taken",
			});
			return;
		}

		var newUser = db.CreateUserParams{
			Email: RegisterData.Email,
			Username: RegisterData.Username,
			PasswordHash: string(hashedPassword),
			PasswordSalt: saltHex,
			FirstName: pgtype.Text{ String: RegisterData.FirstName, Valid: true, },
			LastName: pgtype.Text{ String: RegisterData.LastName, Valid: true, },
			IsActive: pgtype.Bool{ Bool: true, Valid: true, },
			Role: pgtype.Text{ String: "user", Valid: true, },
			PasswordChangedAt: pgtype.Timestamptz{ Time: time.Now(), Valid: true, },
		}
	
		id, err := q.CreateUser(context.Background(), newUser);
		if err != nil && err != pgx.ErrNoRows {
			log.Println(err);
			http.Error(w, err.Error(), http.StatusInternalServerError);
			return;
		}

		if err == pgx.ErrNoRows {
			w.Header().Set("Content-Type", "application/json");
			
			json.NewEncoder(w).Encode(ErrorResponse {
				Message: "Something went wrong creating your account",
			});
			return;
		}

		w.Header().Set("Content-Type", "application/json");
		json.NewEncoder(w).Encode(id);
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

		err = q.LoginEvent(context.Background(), db.LoginEventParams{ 
			ID: user.ID, 
			LastLoginAt: pgtype.Timestamptz{ 
				Time: time.Now(),
				Valid: true,
			},
		});

		if err != nil {
			log.Println(err);
			http.Error(w, err.Error(), http.StatusInternalServerError);
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

