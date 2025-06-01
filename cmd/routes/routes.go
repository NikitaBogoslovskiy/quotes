package routes

import (
	"github.com/NikitaBogoslovskiy/quotes/internal/handlers"
	"github.com/gorilla/mux"
)

type Service struct {
	QuotesHandler handlers.QuotesHandler
}

func NewService(service Service) *Service {
	return &service
}

func (s *Service) LoadRoutes(router *mux.Router) {
	quotes := router.PathPrefix("/quotes").Subrouter()
	quotes.HandleFunc("", s.QuotesHandler.Create).Methods("POST")
	quotes.HandleFunc("", s.QuotesHandler.Get).Methods("GET")
	quotes.HandleFunc("/random", s.QuotesHandler.GetRandom).Methods("GET")
	quotes.HandleFunc("/{id}", s.QuotesHandler.Delete).Methods("DELETE")
}
