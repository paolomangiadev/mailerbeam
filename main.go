// file: main.go
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

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
	err := godotenv.Load(".env", "sendgrid.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	r := Routes()
	http.Handle("/", r)
	log.Printf("Server Listening on %v\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
