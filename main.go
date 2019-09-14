// file: main.go
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	api "github.com/paolomangiadev/mailerbeam/cmd/api"
	auth "github.com/paolomangiadev/mailerbeam/cmd/auth"
)

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Routes definition
func Routes() *chi.Mux {
	router := chi.NewRouter()

	// Basic CORS
	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	login := auth.Init()

	// Router Middlewares
	router.Use(
		cors.Handler,
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.RedirectSlashes,
		middleware.Recoverer,
		jwtauth.Verifier(login),
	)

	// Api Routes
	router.Route("/v1", func(r chi.Router) {
		r.Use(jwtauth.Authenticator)
		r.Mount("/api", api.Routes())
	})

	// Authentication Routes
	router.Mount("/auth", auth.Routes())

	return router
}

type User struct {
	gorm.Model
	Username string
	Email    string `gorm:"type:varchar(100);unique_index"`
	Role     string `gorm:"size:255"`
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db, err := gorm.Open("sqlite3", "mailerbeam.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&User{})
	port := os.Getenv("PORT")
	r := Routes()
	http.Handle("/", r)
	log.Printf("Server Listening on %v\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
