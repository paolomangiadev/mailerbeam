package auth

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	valid "github.com/asaskevich/govalidator"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/paolomangiadev/mailerbeam/app/models"
	renderPkg "github.com/unrolled/render"
)

// Response structure
type Response struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

// TokenResponse structure
type TokenResponse struct {
	Response
	Token string `json:"token,omitempty"`
}

// ValidationErrors type
type ValidationErrors struct {
	Errors map[string]string `json:"validationErrors"`
}

var tokenAuth *jwtauth.JWTAuth
var render *renderPkg.Render

func init() {
	render = renderPkg.New()
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

// Register Handler
func Register(w http.ResponseWriter, req *http.Request) {
	// Read Body
	decoder := json.NewDecoder(req.Body)
	var body models.CreateUserRequest
	err := decoder.Decode(&body)
	if err != nil {
		panic(err)
	}

	// Validate Body
	if isValid, err := valid.ValidateStruct(body); !isValid {
		errors := valid.ErrorsByField(err)
		render.JSON(w, 200, ValidationErrors{errors})
	}

	// Check if user exists
	db := models.GetDB()
	var user models.User
	db.Where(&models.User{
		Email:    body.Email,
		Username: body.Username,
	}).First(&user)
	print(&user)
	// db.Create(&user)
	// render.JSON(w, 200, map[string]interface{}{"user": user})
}

// Login Handler
func Login(w http.ResponseWriter, req *http.Request) {
	tokenClaims := jwt.MapClaims{
		"user_id": 123,
		"exp":     time.Now().Add(time.Hour * time.Duration(4000)).Unix(),
		"iat":     time.Now().Unix(),
		"sub":     123,
	}
	_, tokenString, _ := tokenAuth.Encode(tokenClaims)

	render.JSON(w, 200, TokenResponse{Response{"success", "User logged"}, tokenString})
}
