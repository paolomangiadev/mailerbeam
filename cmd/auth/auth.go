package auth

import (
	"net/http"
	"os"
	"time"

	"encoding/json"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
)

var tokenAuth *jwtauth.JWTAuth

// Response structure
type Response struct {
	Token string `json:"token"`
}

// Routes Auth definitions
func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/login", Login)
	return router
}

// Login Route
func Login(w http.ResponseWriter, req *http.Request) {
	tokenAuth = jwtauth.New("HS256", []byte(os.Getenv("SECRET")), nil)
	tokenClaims := jwt.MapClaims{
		"user_id": 123,
		"exp":     3656,
		"iat":     time.Now().Unix(),
		"sub":     123,
	}
	_, tokenString, _ := tokenAuth.Encode(tokenClaims)

	res, err := json.Marshal(Response{tokenString})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
