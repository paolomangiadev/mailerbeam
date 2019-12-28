package auth

import (
	"encoding/json"
	"fmt"
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
	Token string `json:"token"`
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
	var userbody models.CreateUserRequest
	err := decoder.Decode(&userbody)
	if err != nil {
		panic(err)
	}

	// Validate Body
	if isValid, err := valid.ValidateStruct(userbody); !isValid {
		errors := valid.ErrorsByField(err)
		render.JSON(w, 200, ValidationErrors{errors})
	} else {
		w.Write([]byte(fmt.Sprintf("user registered %v!!!", isValid)))
	}
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

	render.JSON(w, 200, Response{tokenString})
}
