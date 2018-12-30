package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	mux "github.com/gorilla/mux"
	gomail "gopkg.in/gomail.v2"
)

// MailerHandler is the handler function for the index mail route
func MailerHandler(r *mux.Router) {
	r.HandleFunc("/mail", EmailHandler).Methods("POST")
}

// Mail structure
type Mail struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

// EmailHandler controller
func EmailHandler(w http.ResponseWriter, req *http.Request) {

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

	// Email Config
	d := gomail.NewDialer("smtp.gmail.com", 587, "paolo.mangia.dev@gmail.com", "flfyicqttxpszjlf")
	s, err := d.Dial()
	if err != nil {
		panic(err)
	}

	// Sync Email Waitgroup
	var wg sync.WaitGroup
	wg.Add(len(list))

	// New Gomail Message
	m := gomail.NewMessage()

	for _, account := range list {
		go func(s gomail.SendCloser, m *gomail.Message, account Mail) {
			defer wg.Done()
			m.SetHeader("From", "no-reply@example.com")
			m.SetAddressHeader("To", account.Address, account.Name)
			m.SetHeader("Subject", "Newsletter #1")
			m.SetBody("text/html", `<p style="color:red">Mimmoooooo</p>`)

			if err := gomail.Send(s, m); err != nil {
				log.Printf("Could not send email to %q: %v", account.Address, err)
			}
			m.Reset()
		}(s, m, account)
	}

	wg.Wait()

	fmt.Fprintln(w, "Email Sent")
}
