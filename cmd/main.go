package main

import (
	"fmt"
	"net/http"

	"github.com/NikitaBogoslovskiy/quotes/cmd/routes"
	"github.com/NikitaBogoslovskiy/quotes/internal/di"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	service := routes.NewService(routes.Service{
		QuotesHandler: di.InitializeQuotesHandler(),
	})
	service.LoadRoutes(router)

	port := 8080
	http.Handle("/", router)
	fmt.Printf("Start listening to %d port...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
