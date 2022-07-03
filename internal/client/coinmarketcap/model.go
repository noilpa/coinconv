package coinmarketcap

import (
	"encoding/json"
	"fmt"
	"time"
)

type CryptocurrencyListing struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	Symbol      string          `json:"symbol"`
	Slug        string          `json:"slug"`
	LastUpdated time.Time       `json:"last_updated"`
	QuoteRaw    json.RawMessage `json:"quote"`
	Quote       map[string]Quote
}

type Quote struct {
	Price       float64   `json:"price"`
	LastUpdated time.Time `json:"last_updated"`
}

type status struct {
	Timestamp    time.Time `json:"timestamp"`
	ErrorCode    int64     `json:"error_code"`
	ErrorMessage string    `json:"error_message"`
	Elapsed      int64     `json:"elapsed"`
	CreditCount  int64     `json:"credit_count"`
}

func (s status) error() error {
	if s.ErrorCode == 0 && len(s.ErrorMessage) == 0 {
		return nil
	}
	return fmt.Errorf("msg=%s status=%d", s.ErrorMessage, s.ErrorCode)
}

type response struct {
	Status status
	Data   json.RawMessage
}
