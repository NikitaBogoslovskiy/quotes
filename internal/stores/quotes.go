package stores

import (
	"fmt"
	"maps"
	"math"
	"slices"
	"sync"

	"github.com/NikitaBogoslovskiy/quotes/internal/types"
	"golang.org/x/exp/rand"
)

type QuotesStore interface {
	Create(author types.Author, quote types.Quote) (types.Id, error)
	GetAll() ([]types.QuoteData, error)
	GetRandom() (types.QuoteData, error)
	GetByAuthor(author types.Author) (types.QuoteData, error)
	DeleteById(id types.Id) error
}

type quotesStore struct {
	mtx    sync.Mutex
	currId types.Id
	data   map[types.Id]types.QuoteData
}

func NewQuotesStore() QuotesStore {
	return &quotesStore{
		currId: 0,
		data:   make(map[types.Id]types.QuoteData),
	}
}

func (qs *quotesStore) Create(author types.Author, quote types.Quote) (types.Id, error) {
	qs.mtx.Lock()
	defer qs.mtx.Unlock()

	if qs.currId == math.MaxInt64 {
		return 0, fmt.Errorf("Space limit exceeded") // in real-life case it is better to use database and uuid instead of integer ids
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

	return slices.Collect(maps.Values(qs.data)), nil
}

func (qs *quotesStore) GetRandom() (types.QuoteData, error) {
	qs.mtx.Lock()
	defer qs.mtx.Unlock()

	quotesNumber := len(qs.data)
	if quotesNumber == 0 {
		return types.QuoteData{}, fmt.Errorf("No quotes to retrieve")
	}

	randomIdx := rand.Intn(quotesNumber)
	for _, quote := range qs.data {
		if randomIdx == 0 {
			return quote, nil
		}
		randomIdx--
	}

	return types.QuoteData{}, fmt.Errorf("Internal server error")
}

func (qs *quotesStore) GetByAuthor(author types.Author) (types.QuoteData, error) {
	qs.mtx.Lock()
	defer qs.mtx.Unlock()

	for _, quote := range qs.data {
		if quote.Author == author {
			return quote, nil
		}
	}

	return types.QuoteData{}, fmt.Errorf("No quote was found")
}

func (qs *quotesStore) DeleteById(id types.Id) error {
	qs.mtx.Lock()
	defer qs.mtx.Unlock()

	_, ok := qs.data[id]
	if !ok {
		return fmt.Errorf("No quote with specified id")
	}

	delete(qs.data, id) // we can perform delete without upper check but it is better to handle the case when id does not exist
	return nil
}
