// file: main.go
package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	api "github.com/paolomangiadev/mailerbeam/cmd/api"
)

// Routes definition
func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)
	router.Route("/v1", func(r chi.Router) {
		r.Mount("/api", api.Routes())
	})
	return router
}

func main() {
	port := ":5000"
	r := Routes()
	http.Handle("/", r)
	log.Printf("Server Listening on %v\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
