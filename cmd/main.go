package main

import (
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

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
