// Setup and start web server

package main

import (
	"net/http"
	"strings"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)


func main () {
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
	
	router.Route("/api", func (r chi.Router) {
		r.Get("/", func (w http.ResponseWriter, r *http.Request) {
			w.Write([] byte("hello world!"));
		});

	});
	
	// Setup static files for React SPA
	workDir, _ := os.Getwd();
	feDir := http.Dir(filepath.Join(workDir, "web/build/static"));
	FileServer(router, "/static", feDir);

	// Routes for serving React site
	router.Mount("/", ReactRouter());
	
	err := http.ListenAndServe(":3333", router)
	if (err != nil) {
		panic("Could not start web server");
	}
}

func ReactRouter () chi.Router {
	cwd, _ := os.Getwd();
	feDir := http.Dir(filepath.Join(cwd, "web/build"));
	index, _ := os.ReadFile(string(feDir) + "/index.html");

	r := chi.NewRouter();
	
	r.Use(middleware.NoCache);

	r.Get("/*", func (w http.ResponseWriter, r *http.Request) {
		w.Write(index);
	});

	return r;
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context());
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*");
		fs := http.StripPrefix(pathPrefix, http.FileServer(root));
		fs.ServeHTTP(w, r);
	})
}

