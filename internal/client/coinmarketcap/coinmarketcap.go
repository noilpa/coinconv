package coinmarketcap

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	secretHeader = "X-CMC_PRO_API_KEY"
)

var EmptyDataErr = errors.New("got empty data")

type httpCli interface {
	Do(req *http.Request) (*http.Response, error)
}

type client struct {
	c         httpCli
	scheme    string
	host      string
	secretKey string
}

func New() (*client, error) {
	cfg, err := readConfig()
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(cfg.CoinmarketcapHost)
	if err != nil {
		return nil, err
	}
	return &client{
		c: &http.Client{
			Timeout: 10 * time.Second,
		},
		host:      u.Host,
		scheme:    u.Scheme,
		secretKey: cfg.CoinmarketcapSecretKey,
	}, nil
}

func (c *client) CryptocurrencyQuotesLatest(ctx context.Context, from, to string) (CryptocurrencyListing, error) {
	r, err := c.do(ctx, http.MethodGet, "v1/cryptocurrency/quotes/latest", url.Values{"convert": {to}, "symbol": {from}}, nil)
	if err != nil {
		return CryptocurrencyListing{}, err
	}

	var rr map[string]CryptocurrencyListing
	if err := json.Unmarshal(r.Data, &rr); err != nil {
		return CryptocurrencyListing{}, err
	}

	if len(rr) == 0 {
		return CryptocurrencyListing{}, fmt.Errorf("CryptocurrencyQuotesLatest: %w", EmptyDataErr)
	}

	if value, exist := rr[from]; exist {
		var mm map[string]Quote
		if err := json.Unmarshal(rr[from].QuoteRaw, &mm); err != nil {
			return CryptocurrencyListing{}, err
		}
		value.Quote = mm
		rr[from] = value
	}

	return rr[from], nil
}

func (c *client) do(ctx context.Context, method, path string, values url.Values, body []byte) (response, error) {
	u := url.URL{
		Scheme:   c.scheme,
		Host:     c.host,
		Path:     path,
		RawQuery: values.Encode(),
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), bytes.NewReader(body))
	if err != nil {
		return response{}, err
	}

	req.Header = http.Header{
		secretHeader: {c.secretKey},
		"Accept":     {"application/json"},
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return response{}, err
	}

	defer resp.Body.Close()
	var rr response
	if err := json.NewDecoder(resp.Body).Decode(&rr); err != nil {
		return response{}, err
	}

	return rr, nil
}
