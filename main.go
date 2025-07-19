// Setup and start web server

package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	Routers "rtrade/routers"
	Config "rtrade/config"
)


func main () {
	var config * Config.Config;
	var err error;

	config, err = Config.Load();
	if err != nil {
		panic("Failed to load application config");
	}

	// Initial router stack setup, https://github.com/go-chi/chi
	var router = chi.NewRouter();

	// Middleware Foundation
	router.Use(middleware.RequestID);
	router.Use(middleware.RealIP);
	router.Use(middleware.Logger);
	router.Use(middleware.Recoverer);

	// Setup timeout handling
	router.Use(middleware.Timeout(60 * time.Second));

	// Use /api/v1 for all api related queries
	// TODO auth middleware to authenticate requests
	// TODO auth specific router that sets up user login and regi
	// router.Route("/auth", func (r chi.Router) {
	// 	r.Get(...	
	// });
		

	// router.Get("/", func (w http.ResponseWriter, r *http.Request) {
	// 	w.Write([] byte("hello world!"));
	// });
	
	// Handle API
	router.Mount("/api", Routers.Api());
	
	// Setup static files for React SPA
	router.Mount("/static", Routers.Static());

	// Routes for serving React site
	router.Mount("/", Routers.React());
	
	err = http.ListenAndServe(":" + config.Application.Port, router)
	if (err != nil) {
		panic("Could not start web server");
	}
}

