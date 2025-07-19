package routers

import (
	"net/http"
	"path/filepath"
	"strings"
	"os"
	
	"github.com/go-chi/chi/v5"
)

func Static () chi.Router {
	r := chi.NewRouter();

	workDir, _ := os.Getwd();
	feDir := http.Dir(filepath.Join(workDir, "/web/build/static"));
	path := "/";

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
		fs := http.StripPrefix(pathPrefix, http.FileServer(feDir));
		fs.ServeHTTP(w, r);
	})

	return r;
}
