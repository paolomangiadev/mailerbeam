// file: main.go
package main

import (
	"log"
	"net/http"

	routes "mailerbeam/cmd/routes"
)

func main() {
	port := ":5000"
	r := routes.NewRouter()
	http.Handle("/", r)
	log.Printf("Server Listening on %v\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
