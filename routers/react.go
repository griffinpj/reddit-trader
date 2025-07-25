package routers

import (
	"log"
	"net/http"
	"os"
	"rtrade/config"

	// "path/filepath"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/olivere/vite"
)

func React (env * config.Env) chi.Router {
	// cwd, _ := os.Getwd();
	// feDir := http.Dir(filepath.Join(cwd, "/web/build"));
	// index, _ := os.ReadFile(string(feDir) + "/index.html");
	viteFragment, err := vite.HTMLFragment(vite.Config{
		FS: os.DirFS("web"),
		IsDev: true,
		ViteURL: "http://localhost:5173/",
		ViteEntry: "src/index.jsx",
		ViteTemplate: vite.React,
	})
	if err != nil {
		log.Fatal(err);
	}

	indexTemplate := `
		<head>
			<meta charset="UTF-8" />
			<title>My Go Application</title>
			{{ .Vite.Tags }}
		</head>
		<body>
			<div id="root"></div>
		</body>
	`;

	tmpl := template.Must(template.New("name").Parse(indexTemplate));

	r := chi.NewRouter();
	
	r.Use(middleware.NoCache);
	
	r.Group(func (r chi.Router) {
		r.Use(env.Jwt.RequireNoAuth)
		r.HandleFunc("/login", func (w http.ResponseWriter, r *http.Request) {
			pageData := map[string]interface{}{
				"Vite": viteFragment,
			}
			
			if err = tmpl.Execute(w, pageData); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError);
			}
		});
	});

	r.Group(func (r chi.Router) {
		r.Use(env.Jwt.RequireAuthRedir)
		r.HandleFunc("/*", func (w http.ResponseWriter, r *http.Request) {
			pageData := map[string]interface{}{
				"Vite": viteFragment,
			}
			
			if err = tmpl.Execute(w, pageData); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError);
			}
		});
	});


	return r;
}

