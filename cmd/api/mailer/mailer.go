package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/chi"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Routes is the handler function for the index mail route
func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/", SendEmail)
	return router
}

// Mail structure
type Mail struct {
	Name    string `json:"name"`
	Address string `json:"address"`
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

// SendEmail controller
func SendEmail(w http.ResponseWriter, req *http.Request) {

	// Read Body
	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// The list of newsletter emails.
	var list []Mail

	// Unmarshal the body
	err = json.Unmarshal([]byte(body), &list)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Sync Email Waitgroup
	var wg sync.WaitGroup
	wg.Add(len(list))

	plainTextContent := "blast emails to your contacts"
	htmlContent := "<strong>blast emails to your contacts</strong>"
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))

	for _, account := range list {
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
