package types

import "fmt"

type Id uint64

func (id Id) Validate() error {
	if id == 0 {
		return fmt.Errorf("id cannot be zero")
	}

	return nil
}

type Author string

func (a Author) Validate() error {
	if len(a) == 0 {
		return fmt.Errorf("author cannot be empty")
	}

	return nil
}

type Quote string

func (q Quote) Validate() error {
	if len(q) == 0 {
		return fmt.Errorf("quote cannot be empty")
	}

	return nil
}

type QuoteData struct {
	Id     Id     `json:"id,omitempty"`
	Author Author `json:"author,omitempty"`
	Quote  Quote  `json:"quote,omitempty"`
}

type CreateQuoteRequest struct {
	Author Author `json:"author"`
	Quote  Quote  `json:"quote"`
}

func (cqr CreateQuoteRequest) Validate() error {
	err := cqr.Author.Validate()
	if err != nil {
		return err
	}

	err = cqr.Quote.Validate()
	if err != nil {
		return err
	}

	return nil
}

type CreateQuoteResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
	Id      Id     `json:"id,omitempty"`
}

type GetQuotesResponse struct {
	Ok      bool        `json:"ok"`
	Message string      `json:"message,omitempty"`
	Quotes  []QuoteData `json:"quotes"`
}

type GetRandomQuoteResponse struct {
	Ok      bool      `json:"ok"`
	Message string    `json:"message,omitempty"`
	Quote   QuoteData `json:"quote,omitempty"`
}

type DeleteQuoteResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
}
