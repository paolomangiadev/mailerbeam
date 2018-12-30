package api

import (
	mailer "mailerbeam/cmd/api/mailer"

	"github.com/gorilla/mux"
)

// Router of Api
func Router() *mux.Router {
	r := mux.NewRouter()
	mailer.MailerHandler(r)
	return r
}
