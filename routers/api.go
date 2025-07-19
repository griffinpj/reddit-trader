package routers

import (
	"net/http"
	
	"github.com/go-chi/chi/v5"

	Lib "rtrade/lib"
)

func Api () chi.Router {
	r := chi.NewRouter();
	
	r.Get("/", func (w http.ResponseWriter, r *http.Request) {
		Lib.Database();
		w.Write([] byte("hello world!"));
	});

	return r;
}

