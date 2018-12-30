package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi"

	gomail "gopkg.in/gomail.v2"
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
func (mail *Mail) Send(s gomail.SendCloser, wg *sync.WaitGroup, w http.ResponseWriter) {
	defer wg.Done()
	// New Gomail Message
	m := gomail.NewMessage()
	m.SetHeader("From", "no-reply@example.com")
	m.SetAddressHeader("To", mail.Address, mail.Name)
	m.SetHeader("Subject", "Newsletter #1")
	m.SetBody("text/html", `<p style="color:red">New Message</p>`)
	fmt.Printf("%v \n", s)
	if err := gomail.Send(s, m); err != nil {
		log.Print(err)
		fmt.Fprintln(w, "Can't send Email")
	}
	m.Reset()
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

	// Email Config
	d := gomail.NewDialer(
		"smtp.gmail.com",
		587,
		"paolo.mangia.dev@gmail.com",
		"flfyicqttxpszjlf",
	)
	s, err := d.Dial()
	fmt.Printf("%v \n", s)
	if err != nil {
		panic(err)
	}

	// Sync Email Waitgroup
	var wg sync.WaitGroup
	wg.Add(len(list))

	for _, mail := range list {
		acc := mail
		go acc.Send(s, &wg, w)
	}

	wg.Wait()
	s.Close()
	fmt.Fprintln(w, "Email Sent")
}
