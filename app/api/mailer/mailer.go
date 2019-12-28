package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Routes is the handler function for the index mail route
func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/", SendEmail)
	router.Get("/", Protected)
	return router
}

// Mail structure
type Mail struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

// BodyMail structure
type BodyMail struct {
	Contacts []Mail `json:"contacts"`
	HTML     string `json:"html"`
	Title    string `json:"title"`
}

// Send Email Method
func (acc *Mail) Send(wg *sync.WaitGroup, client *sendgrid.Client, message *mail.SGMailV3) {
	defer wg.Done()
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}

// Protected controller
func Protected(w http.ResponseWriter, req *http.Request) {
	_, claims, _ := jwtauth.FromContext(req.Context())
	w.Write([]byte(fmt.Sprintf("PROTECTED AREA. Hi %v!!!", claims["user_id"])))
}

// SendEmail controller
func SendEmail(w http.ResponseWriter, req *http.Request) {

	// Read Body
	decoder := json.NewDecoder(req.Body)
	var mailbody BodyMail
	err := decoder.Decode(&mailbody)
	if err != nil {
		panic(err)
	}

	// Sync Email Waitgroup
	var wg sync.WaitGroup
	wg.Add(len(mailbody.Contacts))

	plainTextContent := "Hello"
	htmlContent := mailbody.HTML
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))

	for _, account := range mailbody.Contacts {
		acc := account
		from := mail.NewEmail("MailerBeam", "no-replay@mailerbeam.com")
		subject := fmt.Sprintf("Welcome %v, start sending mails with MailerBeam", acc.Name)
		to := mail.NewEmail(acc.Name, acc.Address)
		message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
		go acc.Send(&wg, client, message)
	}

	wg.Wait()
	fmt.Fprintln(w, "Email Sent")
}
