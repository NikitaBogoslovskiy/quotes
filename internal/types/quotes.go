package types

import "fmt"

type Id int64

func (id Id) Validate() error {
	if id <= 0 {
		return fmt.Errorf("Id cannot be non-positive")
	}

	return nil
}

type Author string

func (a Author) Validate() error {
	if len(a) == 0 {
		return fmt.Errorf("Author field cannot be empty")
	}

	return nil
}

type Quote string

func (q Quote) Validate() error {
	if len(q) == 0 {
		return fmt.Errorf("Quote field cannot be empty")
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

type GetAllQuotesResponse struct {
	Ok      bool        `json:"ok"`
	Message string      `json:"message,omitempty"`
	Quotes  []QuoteData `json:"quotes,omitempty"`
}

type GetRandomQuoteResponse struct {
	Ok      bool      `json:"ok"`
	Message string    `json:"message,omitempty"`
	Quote   QuoteData `json:"quote,omitempty"`
}

type GetQuoteByAuthorResponse struct {
	Ok      bool      `json:"ok"`
	Message string    `json:"message,omitempty"`
	Quote   QuoteData `json:"quote,omitempty"`
}

type DeleteQuoteByIdResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
}
