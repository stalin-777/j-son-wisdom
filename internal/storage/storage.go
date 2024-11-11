package storage

import (
	_ "embed"
	"encoding/json"
	"math/rand"
	"time"
)

const jason string = " © Джейсон Стэтхэм"

//go:embed quotes.json
var quotesByte []byte

// Storage represents a collection of quotes.
type Storage struct {
	quotes []string
}

// New creates a new instance of the quote store
func New() (*Storage, error) {
	s := &Storage{}

	err := json.Unmarshal(quotesByte, &s.quotes)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// GetRandomQuote retrieves and returns a random quote from the storage.
func (s *Storage) GetRandomQuote() string {
	rand.New(rand.NewSource(time.Now().Unix()))

	return s.quotes[rand.Intn(len(s.quotes))] + jason
}
