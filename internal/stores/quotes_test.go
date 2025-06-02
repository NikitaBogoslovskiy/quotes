package stores

import (
	"fmt"
	"slices"
	"testing"

	"github.com/NikitaBogoslovskiy/quotes/internal/types"
)

func TestCreate(t *testing.T) {
	quotesStore := &quotesStore{data: make(map[types.Id]types.QuoteData)}

	type input struct {
		Author types.Author
		Quote  types.Quote
	}

	type output struct {
		Id  types.Id
		Err error
	}

	testCases := []struct {
		name     string
		input    input
		expected output
	}{
		{
			name:     "FirstRequest",
			input:    input{Author: "Author1", Quote: "Quote1"},
			expected: output{Id: 1, Err: nil},
		},
		{
			name:     "SecondRequest",
			input:    input{Author: "Author2", Quote: "Quote2"},
			expected: output{Id: 2, Err: nil},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id, err := quotesStore.Create(tc.input.Author, tc.input.Quote)
			if err != tc.expected.Err {
				t.Errorf("store returned unexpected error: got %v want %v", err, tc.expected.Err)
			}
			if id != tc.expected.Id {
				t.Errorf("store returned unexpected id: got %v want %v", id, tc.expected.Id)
			}

			if err == nil {
				record, ok := quotesStore.data[id]
				if !ok {
					t.Errorf("store did not save record: no quote with id = %v", id)
				} else {
					if record.Author != tc.input.Author {
						t.Errorf("store saved record with wrong author: got %v want %v", record.Author, tc.input.Author)
					}
					if record.Quote != tc.input.Quote {
						t.Errorf("store saved record with wrong quote: got %v want %v", record.Quote, tc.input.Quote)
					}
				}
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	quotesStore := &quotesStore{data: make(map[types.Id]types.QuoteData)}

	expectedQuotes := make([]types.QuoteData, 0)
	t.Run("NoQuotes", func(t *testing.T) {
		quotes, err := quotesStore.GetAll()
		if err != nil {
			t.Errorf("store returned unexpected error: %v", err)
		}
		if !slices.Equal(quotes, expectedQuotes) {
			t.Errorf("store returned unexpected quotes: got %v want %v", quotes, expectedQuotes)
		}
	})

	var (
		author1 types.Author = "Author1"
		quote1  types.Quote  = "Quote1"
		author2 types.Author = "Author2"
		quote2  types.Quote  = "Quote2"
	)
	id1, _ := quotesStore.Create(author1, quote1)
	id2, _ := quotesStore.Create(author2, quote2)
	expectedQuotes = []types.QuoteData{
		{
			Id:     id1,
			Author: author1,
			Quote:  quote1,
		},
		{
			Id:     id2,
			Author: author2,
			Quote:  quote2,
		},
	}
	t.Run("TwoQuotes", func(t *testing.T) {
		quotes, err := quotesStore.GetAll()
		if err != nil {
			t.Errorf("store returned unexpected error: %v", err)
		}
		slices.SortFunc(quotes, func(q1, q2 types.QuoteData) int {
			return int(q1.Id - q2.Id)
		})
		if !slices.Equal(quotes, expectedQuotes) {
			t.Errorf("store returned unexpected quotes: got %v want %v", quotes, expectedQuotes)
		}
	})
}

func TestGetByAuthor(t *testing.T) {
	quotesStore := &quotesStore{data: make(map[types.Id]types.QuoteData)}

	var (
		author1 types.Author = "Author1"
		author2 types.Author = "Author2"
		quote1  types.Quote  = "Quote1"
		quote2  types.Quote  = "Quote2"
		quote3  types.Quote  = "Quote3"
	)
	id1, _ := quotesStore.Create(author1, quote1)
	quotesStore.Create(author2, quote2)
	id3, _ := quotesStore.Create(author1, quote3)

	expectedQuotes := make([]types.QuoteData, 0)
	t.Run("WrongAuthor", func(t *testing.T) {
		quotes, err := quotesStore.GetByAuthor("abcd")
		if err != nil {
			t.Errorf("store returned unexpected error: %v", err)
		}
		if !slices.Equal(quotes, expectedQuotes) {
			t.Errorf("store returned unexpected quotes: got %v want %v", quotes, expectedQuotes)
		}
	})

	expectedQuotes = []types.QuoteData{
		{
			Id:     id1,
			Author: author1,
			Quote:  quote1,
		},
		{
			Id:     id3,
			Author: author1,
			Quote:  quote3,
		},
	}
	t.Run("CorrectAuthor", func(t *testing.T) {
		quotes, err := quotesStore.GetByAuthor(author1)
		if err != nil {
			t.Errorf("store returned unexpected error: %v", err)
		}
		slices.SortFunc(quotes, func(q1, q2 types.QuoteData) int {
			return int(q1.Id - q2.Id)
		})
		if !slices.Equal(quotes, expectedQuotes) {
			t.Errorf("store returned unexpected quotes: got %v want %v", quotes, expectedQuotes)
		}
	})
}

func TestGetRandom(t *testing.T) {
	quotesStore := &quotesStore{data: make(map[types.Id]types.QuoteData)}

	expectedError := fmt.Errorf("no quotes to retrieve")
	t.Run("EmptyStore", func(t *testing.T) {
		_, err := quotesStore.GetRandom()
		if err == nil || err.Error() != expectedError.Error() {
			t.Errorf("store returned unexpected error: got %v want %v", err, expectedError)
		}
	})

	var (
		author1 types.Author = "Author1"
		quote1  types.Quote  = "Quote1"
		author2 types.Author = "Author2"
		quote2  types.Quote  = "Quote2"
	)
	id1, _ := quotesStore.Create(author1, quote1)
	id2, _ := quotesStore.Create(author2, quote2)
	expectedQuoteIds := map[types.Id]bool{id1: true, id2: true}
	t.Run("StoreWithTwoQuotes", func(t *testing.T) {
		quote, err := quotesStore.GetRandom()
		if err != nil {
			t.Errorf("store returned unexpected error: %v", err)
		}
		if !expectedQuoteIds[quote.Id] {
			t.Errorf("store returned quote with unexpected id: %v", quote.Id)
		}
	})
}

func TestDelete(t *testing.T) {
	quotesStore := &quotesStore{data: make(map[types.Id]types.QuoteData)}

	expectedError := fmt.Errorf("no quote with specified id")
	t.Run("EmptyStore", func(t *testing.T) {
		err := quotesStore.Delete(types.Id(1))
		if err == nil || err.Error() != expectedError.Error() {
			t.Errorf("store returned unexpected error: got %v want %v", err, expectedError)
		}
	})

	var (
		author1 types.Author = "Author1"
		quote1  types.Quote  = "Quote1"
		author2 types.Author = "Author2"
		quote2  types.Quote  = "Quote2"
	)
	id1, _ := quotesStore.Create(author1, quote1)
	quotesStore.Create(author2, quote2)

	t.Run("WrongId", func(t *testing.T) {
		err := quotesStore.Delete(types.Id(50))
		if err == nil || err.Error() != expectedError.Error() {
			t.Errorf("store returned unexpected error: got %v want %v", err, expectedError)
		}
	})

	t.Run("CorrectId", func(t *testing.T) {
		err := quotesStore.Delete(id1)
		if err != nil {
			t.Errorf("store returned unexpected error: %v", err)
		}

		_, ok := quotesStore.data[id1]
		if ok {
			t.Errorf("store did not delete quote with id = %v", id1)
		}
	})
}
