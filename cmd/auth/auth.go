package auth

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"encoding/json"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/jinzhu/gorm"
	models "github.com/paolomangiadev/mailerbeam/app/models"
)

var tokenAuth *jwtauth.JWTAuth

// Response structure
type Response struct {
	Token string `json:"token"`
}

type BodyRegister struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Routes Auth definitions
func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/login", Login)
	router.Post("/register", Register)
	return router
}

// Init Auth
func Init() *jwtauth.JWTAuth {
	tokenAuth = jwtauth.New("HS256", []byte(os.Getenv("SECRET")), nil)
	return tokenAuth
}

// Login Route
func Login(w http.ResponseWriter, req *http.Request) {
	tokenClaims := jwt.MapClaims{
		"user_id": 123,
		"exp":     time.Now().Add(time.Hour * time.Duration(4000)).Unix(),
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

func Register(w http.ResponseWriter, req *http.Request) {
	regbody := BodyRegister{}
	// Unmarshal the body
	err := json.NewDecoder(req.Body).Decode(&regbody)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	db := models.GetDB()
	var user models.User
	if err := db.Where("email = ?", regbody.Email).First(&user).Error; gorm.IsRecordNotFoundError(err) {
		// record not found
		w.Write([]byte(fmt.Sprintf("User with email: %v not found!!!", regbody.Email)))
	}
}
