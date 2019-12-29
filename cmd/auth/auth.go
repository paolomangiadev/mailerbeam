package auth

import (
	"encoding/json"
	"log"
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

// DataResponse structure
type DataResponse struct {
	Response
	Data models.User `json:"data"`
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
		log.Println(errors)
		render.JSON(w, 500, ValidationErrors{errors})
		return
	}

	// Check if user exists
	db := models.GetDB()
	var user models.User
	if db.Where(&models.User{
		Email:    body.Email,
		Username: body.Username,
	}).First(&user).RecordNotFound() {
		user := models.User{
			Name:     body.Name,
			Email:    body.Email,
			Username: body.Username,
			Password: body.Password,
			Role:     "admin",
		}
		// Create user if doesn't exists
		if err := db.Create(&user).Error; err != nil {
			render.JSON(w, 500, Response{"fail", "Unable to create user, error: " + err.Error()})
			return
		}
		render.JSON(w, 200, DataResponse{Response{"success", "User Created"}, user})
		return
	}
	render.JSON(w, 403, Response{"fail", "User already exists"})
	return
}

// Login Handler
func Login(w http.ResponseWriter, req *http.Request) {

	// Read Body
	decoder := json.NewDecoder(req.Body)
	var body models.LogingUserRequest
	err := decoder.Decode(&body)
	if err != nil {
		panic(err)
	}

	// Validate Body
	if isValid, err := valid.ValidateStruct(body); !isValid {
		errors := valid.ErrorsByField(err)
		log.Println(errors)
		render.JSON(w, 500, ValidationErrors{errors})
		return
	}

	// Check if user exists
	db := models.GetDB()
	var user models.User
	if db.Where(&models.User{
		Email: body.Email,
	}).First(&user).RecordNotFound() {
		render.JSON(w, 200, Response{"fail", "User doesn't exists"})
		return
	}

	tokenClaims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * time.Duration(4000)).Unix(),
		"iat":     time.Now().Unix(),
		"sub":     123,
	}
	_, tokenString, _ := tokenAuth.Encode(tokenClaims)

	render.JSON(w, 200, TokenResponse{Response{"success", "User logged"}, tokenString})
}
