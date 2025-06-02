package services

import (
	"slices"
	"testing"

	"github.com/NikitaBogoslovskiy/quotes/internal/types"
)

type quotesStoreStub struct{}

func (qs *quotesStoreStub) Create(author types.Author, quote types.Quote) (types.Id, error) {
	return 1, nil
}

func (qs *quotesStoreStub) GetAll() ([]types.QuoteData, error) {
	return make([]types.QuoteData, 2), nil
}

func (qs *quotesStoreStub) GetByAuthor(author types.Author) ([]types.QuoteData, error) {
	return make([]types.QuoteData, 1), nil
}

func (qs *quotesStoreStub) GetRandom() (types.QuoteData, error) {
	return types.QuoteData{}, nil
}

func (qs *quotesStoreStub) Delete(id types.Id) error {
	return nil
}

func TestCreate(t *testing.T) {
	quotesService := NewQuotesService(&quotesStoreStub{})

	testCases := []struct {
		name     string
		input    types.CreateQuoteRequest
		expected types.CreateQuoteResponse
	}{
		{
			name:     "EmptyRequest",
			input:    types.CreateQuoteRequest{},
			expected: types.CreateQuoteResponse{Ok: false, Message: "author cannot be empty"},
		},
		{
			name:     "EmptyAuthor",
			input:    types.CreateQuoteRequest{Quote: "Quote"},
			expected: types.CreateQuoteResponse{Ok: false, Message: "author cannot be empty"},
		},
		{
			name:     "EmptyQuote",
			input:    types.CreateQuoteRequest{Author: "Author"},
			expected: types.CreateQuoteResponse{Ok: false, Message: "quote cannot be empty"},
		},
		{
			name:     "CorrectRequest",
			input:    types.CreateQuoteRequest{Author: "Author", Quote: "Quote"},
			expected: types.CreateQuoteResponse{Ok: true, Id: 1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := quotesService.Create(tc.input)
			if got != tc.expected {
				t.Errorf("service returned unexpected response: got %v want %v",
					got, tc.expected)
			}
		})
	}
}

func TestGet(t *testing.T) {
	quotesService := NewQuotesService(&quotesStoreStub{})

	testCases := []struct {
		name     string
		input    types.Author
		expected types.GetQuotesResponse
	}{
		{
			name:     "EmptyAuthor",
			input:    types.Author(""),
			expected: types.GetQuotesResponse{Ok: true, Quotes: make([]types.QuoteData, 2)},
		},
		{
			name:     "SpecifiedAuthor",
			input:    types.Author("Author"),
			expected: types.GetQuotesResponse{Ok: true, Quotes: make([]types.QuoteData, 1)},
		},
	}

	responsesEqual := func(got, expected types.GetQuotesResponse) bool {
		if got.Ok != expected.Ok {
			return false
		}

		if got.Message != expected.Message {
			return false
		}

		if !slices.Equal(got.Quotes, expected.Quotes) {
			return false
		}

		return true
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := quotesService.Get(tc.input)
			if !responsesEqual(got, tc.expected) {
				t.Errorf("service returned unexpected response: got %v want %v",
					got, tc.expected)
			}
		})
	}
}

func TestGetRandom(t *testing.T) {
	quotesService := NewQuotesService(&quotesStoreStub{})

	testCases := []struct {
		name     string
		expected types.GetRandomQuoteResponse
	}{
		{
			name:     "GetRandom",
			expected: types.GetRandomQuoteResponse{Ok: true, Quote: types.QuoteData{}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := quotesService.GetRandom()
			if got != tc.expected {
				t.Errorf("service returned unexpected response: got %v want %v",
					got, tc.expected)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	quotesService := NewQuotesService(&quotesStoreStub{})

	testCases := []struct {
		name     string
		input    types.Id
		expected types.DeleteQuoteResponse
	}{
		{
			name:     "ZeroId",
			input:    0,
			expected: types.DeleteQuoteResponse{Ok: false, Message: "id cannot be zero"},
		},
		{
			name:     "CorrectId",
			input:    2,
			expected: types.DeleteQuoteResponse{Ok: true},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := quotesService.Delete(tc.input)
			if got != tc.expected {
				t.Errorf("service returned unexpected response: got %v want %v",
					got, tc.expected)
			}
		})
	}
}
