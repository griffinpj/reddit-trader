package routers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func React () chi.Router {
	cwd, _ := os.Getwd();
	feDir := http.Dir(filepath.Join(cwd, "/web/build"));
	index, _ := os.ReadFile(string(feDir) + "/index.html");

	r := chi.NewRouter();
	
	r.Use(middleware.NoCache);

	r.Get("/*", func (w http.ResponseWriter, r *http.Request) {
		w.Write(index);
	});

	return r;
}

