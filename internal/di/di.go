package di

import (
	"github.com/NikitaBogoslovskiy/quotes/internal/handlers"
	"github.com/NikitaBogoslovskiy/quotes/internal/services"
	"github.com/NikitaBogoslovskiy/quotes/internal/stores"
)

func InitializeQuotesHandler() handlers.QuotesHandler {
	quotesStore := stores.NewQuotesStore()
	quotesService := services.NewQuotesService(quotesStore)
	quotesHandler := handlers.NewQuotesHandler(quotesService)

	return quotesHandler
}
