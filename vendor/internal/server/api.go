package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Start(lines []map[string]string) {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	})
	http.ListenAndServe(":3000", r)
}
