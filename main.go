// Setup and start web server

package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"rtrade/auth"
	Config "rtrade/config"
	Lib "rtrade/lib"
	Routers "rtrade/routers"
)

func main() {
	var config *Config.Config
	var err error

	config, err = Config.Load()
	if err != nil {
		panic("Failed to load application config")
	}



	// Instantiate our DB pool and store in our App Env
	pool := Lib.Database()
	jwt := auth.NewJWTManager(auth.Config{
		SecretKey:      "your-secret-key-here", // Use environment variable in production
		TokenExpiry:    24 * time.Hour,
		CookieName:     "auth_token",
		CookieDomain:   "",    // Set your domain
		CookieSecure:   false, // Set to true in production with HTTPS
		CookieSameSite: http.SameSiteLaxMode,
		Pool: pool,
	})

	env := &Config.Env{
		Pool: pool,
		Jwt: jwt,
		Config: config,
	}

	// Close the db connection once things have finished
	defer pool.Close()

	// Initial router stack setup, https://github.com/go-chi/chi
	var router = chi.NewRouter()

	// Middleware Foundation
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Setup timeout handling
	router.Use(middleware.Timeout(60 * time.Second))

	// Use /api/v1 for all api related queries
	// TODO auth middleware to authenticate requests
	// TODO auth specific router that sets up user login and regi
	// router.Route("/auth", func (r chi.Router) {
	// 	r.Get(...
	// });

	// router.Get("/", func (w http.ResponseWriter, r *http.Request) {
	// 	w.Write([] byte("hello world!"));
	// });

	router.Mount("/auth", Routers.Auth(env))

	// Handle API
	router.Mount("/api/v1", Routers.Api(env))

	// Setup static files for React SPA
	router.Mount("/static", Routers.Static())

	// Routes for serving React site
	router.Mount("/", Routers.React(env))

	err = http.ListenAndServe(":"+config.Application.Port, router)
	if err != nil {
		panic("Could not start web server")
	}
}
