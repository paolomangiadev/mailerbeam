package routes

import (
	"fmt"
	"net/http"
	"strings"

	api "github.com/paolomangiadev/mailerbeam/cmd/api"

	"github.com/gorilla/mux"
)

// NewRouter func
func NewRouter() *mux.Router {
	r := mux.NewRouter()
	mount(r, "/api", api.Router())
	logRoutes(r)
	return r
}

// Route Mounter
func mount(r *mux.Router, path string, handler http.Handler) {
	r.PathPrefix(path).Handler(
		http.StripPrefix(
			strings.TrimSuffix(path, "/"),
			handler,
		),
	)
}

func logRoutes(r *mux.Router) {
	err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			fmt.Println("ROUTE:", pathTemplate)
		}
		methods, err := route.GetMethods()
		if err == nil {
			fmt.Println("Methods:", strings.Join(methods, ","))
		}
		fmt.Println()
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
}
