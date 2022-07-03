package conventer

import (
	"context"
	"errors"
	"fmt"

	"coinconv/internal/client/coinmarketcap"
)

type listingGetter interface {
	CryptocurrencyQuotesLatest(ctx context.Context, from, to string) (coinmarketcap.CryptocurrencyListing, error)
}

type converter struct {
	lg listingGetter
}

func New(lg listingGetter) *converter {
	return &converter{
		lg: lg,
	}
}

func (c *converter) Convert(ctx context.Context, num float64, from, to string) (float64, error) {
	listing, err := c.lg.CryptocurrencyQuotesLatest(ctx, from, to)
	if err == nil {
		if _, exist := listing.Quote[to]; exist {
			return num * listing.Quote[to].Price, nil
		}
	}

	// try to find quotes for vice versa
	if errors.As(err, &coinmarketcap.EmptyDataErr) {
		listing, err = c.lg.CryptocurrencyQuotesLatest(ctx, to, from)
		if err != nil {
			return 0, err
		}

		if _, exist := listing.Quote[from]; exist {
			return num * (1 / listing.Quote[from].Price), nil
		}
	}

	if err != nil {
		return 0, err
	}

	return 0, fmt.Errorf("listing from %s to %s is not found", from, to)
}
