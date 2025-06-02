package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/NikitaBogoslovskiy/quotes/internal/types"
	"github.com/gorilla/mux"
)

type quotesServiceStub struct{}

func (qs *quotesServiceStub) Create(request types.CreateQuoteRequest) types.CreateQuoteResponse {
	return types.CreateQuoteResponse{Ok: true, Id: 1}
}

func (qs *quotesServiceStub) Get(author types.Author) types.GetQuotesResponse {
	return types.GetQuotesResponse{Ok: true, Quotes: make([]types.QuoteData, 0)}
}

func (qs *quotesServiceStub) GetRandom() types.GetRandomQuoteResponse {
	return types.GetRandomQuoteResponse{Ok: true, Quote: types.QuoteData{}}
}

func (qs *quotesServiceStub) Delete(id types.Id) types.DeleteQuoteResponse {
	return types.DeleteQuoteResponse{Ok: true}
}

func TestCreate(t *testing.T) {
	quotesHandler := NewQuotesHandler(&quotesServiceStub{})

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "EmptyInput",
			input:    ``,
			expected: `{"ok":false,"message":"incorrect request format"}`,
		},
		{
			name:     "IncorrectBrackets",
			input:    `{`,
			expected: `{"ok":false,"message":"incorrect request format"}`,
		},
		{
			name:     "IncorrectQuotes",
			input:    `{"author":"Author,"quote":"Quote"}`,
			expected: `{"ok":false,"message":"incorrect request format"}`,
		},
		{
			name:     "IncorrectDataType",
			input:    `{"author":"Author","quote":42}`,
			expected: `{"ok":false,"message":"incorrect request format"}`,
		},
		{
			name:     "CorrectInput",
			input:    `{"author":"Author","quote":"Quote"}`,
			expected: `{"ok":true,"id":1}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/quotes", bytes.NewBuffer([]byte(tc.input)))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(quotesHandler.Create)
			handler.ServeHTTP(rr, req)

			if rr.Body.String() != tc.expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tc.expected)
			}
		})
	}
}

func TestGet(t *testing.T) {
	quotesHandler := NewQuotesHandler(&quotesServiceStub{})

	testCases := []struct {
		name      string
		urlParams map[string]string
		expected  string
	}{
		{
			name:      "NoURLParams",
			urlParams: map[string]string{},
			expected:  `{"ok":true,"quotes":[]}`,
		},
		{
			name:      "UnknownURLParams",
			urlParams: map[string]string{"hello": "world", "field": "42"},
			expected:  `{"ok":true,"quotes":[]}`,
		},
		{
			name:      "AuthorURLParam",
			urlParams: map[string]string{"author": "Author"},
			expected:  `{"ok":true,"quotes":[]}`,
		},
		{
			name:      "AuthorURLParams",
			urlParams: map[string]string{"author": "Author", "other": "field"},
			expected:  `{"ok":true,"quotes":[]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var url string
			urlParamsCount := len(tc.urlParams)
			if urlParamsCount == 0 {
				url = "/quotes"
			} else {
				urlParams := make([]string, 0, urlParamsCount)
				for key, value := range tc.urlParams {
					urlParams = append(urlParams, fmt.Sprintf("%s=%s", key, value))
				}
				url = fmt.Sprintf("/quotes?%s", strings.Join(urlParams, "&"))
			}

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(quotesHandler.Get)
			handler.ServeHTTP(rr, req)

			if rr.Body.String() != tc.expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tc.expected)
			}
		})
	}
}

func TestGetRandom(t *testing.T) {
	quotesHandler := NewQuotesHandler(&quotesServiceStub{})

	testCases := []struct {
		name     string
		expected string
	}{
		{
			name:     "GetRandom",
			expected: `{"ok":true,"quote":{}}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/quotes/random", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(quotesHandler.GetRandom)
			handler.ServeHTTP(rr, req)

			if rr.Body.String() != tc.expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tc.expected)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	quotesHandler := NewQuotesHandler(&quotesServiceStub{})

	testCases := []struct {
		name     string
		quoteId  string
		expected string
	}{
		{
			name:     "EmptyId",
			quoteId:  "",
			expected: `{"ok":false,"message":"id should be a non-negative number"}`,
		},
		{
			name:     "StringId",
			quoteId:  "abc",
			expected: `{"ok":false,"message":"id should be a non-negative number"}`,
		},
		{
			name:     "NegativeId",
			quoteId:  "-2",
			expected: `{"ok":false,"message":"id should be a non-negative number"}`,
		},
		{
			name:     "CorrectId",
			quoteId:  "1",
			expected: `{"ok":true}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("DELETE", "/quotes", nil)
			if err != nil {
				t.Fatal(err)
			}
			req = mux.SetURLVars(req, map[string]string{"id": tc.quoteId})

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(quotesHandler.Delete)
			handler.ServeHTTP(rr, req)

			if rr.Body.String() != tc.expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tc.expected)
			}
		})
	}
}
