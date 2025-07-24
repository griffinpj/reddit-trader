package routers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"rtrade/config"
	"rtrade/db"

	"github.com/go-chi/chi/v5"
)

func Api (env *config.Env) chi.Router {
	r := chi.NewRouter();
	
	r.Get("/users", func (w http.ResponseWriter, r *http.Request) {
		conn, err := env.Pool.Acquire(context.Background())
		if err != nil {
			log.Panic("Error aquiring connection");
		}

		// Return the aquired db connection back to the connection pool from which it came!
		defer conn.Release();

		q := db.New(conn);
		
		var users [] db.User
		users, err = q.GetUsers(r.Context())	
		if err != nil {
			log.Panic("Error getting users");
		}

		log.Println(users);
		w.Header().Set("Content-Type", "application/json");
		
		json.NewEncoder(w).Encode(users);
	});

	return r;
}

