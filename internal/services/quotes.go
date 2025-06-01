package services

import (
	"github.com/NikitaBogoslovskiy/quotes/internal/stores"
	"github.com/NikitaBogoslovskiy/quotes/internal/types"
)

type QuotesService interface {
	Create(request types.CreateQuoteRequest) types.CreateQuoteResponse
	Get(author types.Author) types.GetQuotesResponse
	GetRandom() types.GetRandomQuoteResponse
	Delete(id types.Id) types.DeleteQuoteResponse
}

type quotesService struct {
	quotesStore stores.QuotesStore
}

func NewQuotesService(quotesStore stores.QuotesStore) QuotesService {
	return &quotesService{quotesStore: quotesStore}
}

func (qs *quotesService) Create(request types.CreateQuoteRequest) types.CreateQuoteResponse {
	err := request.Validate()
	if err != nil {
		return types.CreateQuoteResponse{Ok: false, Message: err.Error()}
	}

	id, err := qs.quotesStore.Create(request.Author, request.Quote)
	if err != nil {
		return types.CreateQuoteResponse{Ok: false, Message: err.Error()}
	}

	return types.CreateQuoteResponse{Ok: true, Id: id}
}

func (qs *quotesService) Get(author types.Author) types.GetQuotesResponse {
	var quotes []types.QuoteData

	err := author.Validate()
	if err == nil { // if author param is specified and valid we filter results by author
		quotes, err = qs.quotesStore.GetByAuthor(author)
		if err != nil {
			return types.GetQuotesResponse{Ok: false, Message: err.Error()}
		}
	} else { // otherwise return all results
		quotes, err = qs.quotesStore.GetAll()
		if err != nil {
			return types.GetQuotesResponse{Ok: false, Message: err.Error()}
		}
	}

	return types.GetQuotesResponse{Ok: true, Quotes: quotes}
}

func (qs *quotesService) GetRandom() types.GetRandomQuoteResponse {
	quote, err := qs.quotesStore.GetRandom()
	if err != nil {
		return types.GetRandomQuoteResponse{Ok: false, Message: err.Error()}
	}

	return types.GetRandomQuoteResponse{Ok: true, Quote: quote}
}

func (qs *quotesService) Delete(id types.Id) types.DeleteQuoteResponse {
	err := id.Validate()
	if err != nil {
		return types.DeleteQuoteResponse{Ok: false, Message: err.Error()}
	}

	err = qs.quotesStore.Delete(id)
	if err != nil {
		return types.DeleteQuoteResponse{Ok: false, Message: err.Error()}
	}

	return types.DeleteQuoteResponse{Ok: true}
}
