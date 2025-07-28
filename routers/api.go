package routers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"rtrade/auth"
	Config "rtrade/config"
	"rtrade/db"
	"rtrade/reddit"

	"github.com/go-chi/chi/v5"
)

func Api (env * Config.Env) chi.Router {
	r := chi.NewRouter();
	
	redditClient := reddit.NewClient(&env.Config.Reddit);

	r.Use(env.Jwt.RequireAuth)

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
		
		w.Header().Set("Content-Type", "application/json");
		json.NewEncoder(w).Encode(users);
	});

	r.Get("/reddit/connect", func (w http.ResponseWriter, r *http.Request) {
		var redditUrl = "https://www.reddit.com/api/v1/authorize?" + 
			"client_id=" + env.Config.Reddit.ClientId + 
			"&response_type=code" + 
			"&state=" + env.Config.Reddit.State + 
			"&redirect_uri=" + env.Config.Reddit.RedirectUrl + 
			"&duration=permanent" + 
			"&scope=identity edit flair history read vote wikiread wikiedit";

		http.Redirect(w, r, redditUrl, http.StatusSeeOther);
	});

	r.Post("/reddit/token", func (w http.ResponseWriter, r *http.Request) {
		var ctx = r.Context();

		var TokenData struct {
			Code string `json:"code"`
		}

		decoder := json.NewDecoder(r.Body);
		err := decoder.Decode(&TokenData);
		if err != nil {
			log.Println(err);
			http.Error(w, err.Error(), http.StatusInternalServerError);
			return;
		}

		// Get token
		tokenData, err := redditClient.ExchangeCode(TokenData.Code);
		if err != nil {
			log.Println(err);
			http.Error(w, err.Error(), http.StatusInternalServerError);
			return;
		}
		

		var user = ctx.Value("user").(db.User);
		var claimsData = &auth.ClaimsData{
			UserId: user.ID,
			Email: user.Email,
			Username: user.Username,
			RedditToken: auth.Token{
				AccessToken: tokenData.AccessToken,
				RefreshToken: tokenData.RefreshToken,
				Type: tokenData.TokenType,
			},
		}

		token, err := env.Jwt.GenerateToken(claimsData);
		if err != nil {
			log.Println(err);
			http.Error(w, err.Error(), http.StatusInternalServerError);
			return;
		}
		
		// save updated token with reddit auth token
		env.Jwt.SetAuthCookie(w, token);
		
		w.Header().Set("Content-Type", "application/json");
        json.NewEncoder(w).Encode(map[string]interface{}{
            "access_token": tokenData.AccessToken,
            "token_type":   tokenData.TokenType,
            "expires_in":   tokenData.ExpiresIn,
            "scope":        tokenData.Scope,

            // Be careful with refresh tokens - store them securely!
            "has_refresh_token": tokenData.RefreshToken != "",
        });
	});

	return r;
}

