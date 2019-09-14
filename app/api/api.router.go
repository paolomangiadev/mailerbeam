package api

import (
	"github.com/go-chi/chi"
	mailer "github.com/paolomangiadev/mailerbeam/app/api/mailer"
)

// Routes api definitions
func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Mount("/mail", mailer.Routes())
	return router
}
