package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/NikitaBogoslovskiy/quotes/internal/services"
	"github.com/NikitaBogoslovskiy/quotes/internal/types"
	"github.com/gorilla/mux"
)

type QuotesHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	GetRandom(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type quotesHandler struct {
	quotesService services.QuotesService
}

func NewQuotesHandler(quotesService services.QuotesService) QuotesHandler {
	return &quotesHandler{quotesService: quotesService}
}

func (qh *quotesHandler) Create(w http.ResponseWriter, r *http.Request) {
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		errorRsponse, _ := json.Marshal(types.CreateQuoteResponse{Ok: false, Message: "internal server error"})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorRsponse)
		return
	}

	request := types.CreateQuoteRequest{}
	err = json.Unmarshal(requestBody, &request)
	if err != nil {
		errorRsponse, _ := json.Marshal(types.CreateQuoteResponse{Ok: false, Message: "incorrect request format"})
		w.Write(errorRsponse)
		return
	}

	response := qh.quotesService.Create(request)

	responseBody, err := json.Marshal(response)
	if err != nil {
		errorRsponse, _ := json.Marshal(types.CreateQuoteResponse{Ok: false, Message: "internal server error"})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorRsponse)
		return
	}

	w.Write(responseBody)
}

func (qh *quotesHandler) Get(w http.ResponseWriter, r *http.Request) {
	author := types.Author(r.URL.Query().Get("author"))

	response := qh.quotesService.Get(author)

	responseBody, err := json.Marshal(response)
	if err != nil {
		errorRsponse, _ := json.Marshal(types.GetQuotesResponse{Ok: false, Message: "internal server error"})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorRsponse)
		return
	}

	w.Write(responseBody)
}

func (qh *quotesHandler) GetRandom(w http.ResponseWriter, r *http.Request) {
	response := qh.quotesService.GetRandom()

	responseBody, err := json.Marshal(response)
	if err != nil {
		errorRsponse, _ := json.Marshal(types.GetRandomQuoteResponse{Ok: false, Message: "internal server error"})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorRsponse)
		return
	}

	w.Write(responseBody)
}

func (qh *quotesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		errorRsponse, _ := json.Marshal(types.DeleteQuoteResponse{Ok: false, Message: "id should be a non-negative number"})
		w.Write(errorRsponse)
		return
	}

	response := qh.quotesService.Delete(types.Id(id))

	responseBody, err := json.Marshal(response)
	if err != nil {
		errorRsponse, _ := json.Marshal(types.DeleteQuoteResponse{Ok: false, Message: "internal server error"})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorRsponse)
		return
	}

	w.Write(responseBody)
}
