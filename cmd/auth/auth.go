package auth

import (
	"net/http"

	"github.com/go-chi/chi"
)

// Routes Auth definitions
func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/login", Login)
	return router
}

// Login Route
func Login(w http.ResponseWriter, req *http.Request) {

}
