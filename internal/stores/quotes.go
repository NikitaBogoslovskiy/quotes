package stores

import (
	"fmt"
	"math"
	"math/rand"
	"sync"

	"github.com/NikitaBogoslovskiy/quotes/internal/types"
)

type QuotesStore interface {
	Create(author types.Author, quote types.Quote) (types.Id, error)
	GetAll() ([]types.QuoteData, error)
	GetByAuthor(author types.Author) ([]types.QuoteData, error)
	GetRandom() (types.QuoteData, error)
	Delete(id types.Id) error
}

type quotesStore struct {
	mtx    sync.Mutex
	currId types.Id
	data   map[types.Id]types.QuoteData
}

func NewQuotesStore() QuotesStore {
	return &quotesStore{data: make(map[types.Id]types.QuoteData)}
}

func (qs *quotesStore) Create(author types.Author, quote types.Quote) (types.Id, error) {
	qs.mtx.Lock()
	defer qs.mtx.Unlock()

	if qs.currId == math.MaxUint64 {
		return 0, fmt.Errorf("space limit exceeded")
	}
	qs.currId++

	qs.data[qs.currId] = types.QuoteData{
		Id:     qs.currId,
		Author: author,
		Quote:  quote,
	}

	return qs.currId, nil
}

func (qs *quotesStore) GetAll() ([]types.QuoteData, error) {
	qs.mtx.Lock()
	defer qs.mtx.Unlock()

	quotes := make([]types.QuoteData, 0, len(qs.data))
	for _, quote := range qs.data {
		quotes = append(quotes, quote)
	}

	return quotes, nil
}

func (qs *quotesStore) GetByAuthor(author types.Author) ([]types.QuoteData, error) {
	qs.mtx.Lock()
	defer qs.mtx.Unlock()

	quotes := make([]types.QuoteData, 0)
	for _, quote := range qs.data {
		if quote.Author == author {
			quotes = append(quotes, quote)
		}
	}

	return quotes, nil
}

func (qs *quotesStore) GetRandom() (types.QuoteData, error) {
	qs.mtx.Lock()
	defer qs.mtx.Unlock()

	quotesNumber := len(qs.data)
	if quotesNumber == 0 {
		return types.QuoteData{}, fmt.Errorf("no quotes to retrieve")
	}

	randomIdx := rand.Intn(quotesNumber)
	for _, quote := range qs.data {
		if randomIdx == 0 {
			return quote, nil
		}
		randomIdx--
	}

	return types.QuoteData{}, fmt.Errorf("internal server error")
}

func (qs *quotesStore) Delete(id types.Id) error {
	qs.mtx.Lock()
	defer qs.mtx.Unlock()

	_, ok := qs.data[id]
	if !ok {
		return fmt.Errorf("no quote with specified id")
	}

	delete(qs.data, id)
	return nil
}
