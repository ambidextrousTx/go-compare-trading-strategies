// internal/fetch/alphavantage.go
package fetch

import (
	"fmt"
	"net/http"
	"net/url"

	"go-compare-trading-strategies/internal/model"
)

const baseURL = "https://www.alphavantage.co/query"

// FetchDaily retrieves full daily price history for a ticker.
func FetchDaily(ticker, apiKey string) ([]model.PricePoint, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("parsing base URL: %w", err)
	}

	q := u.Query()
	q.Set("function", "TIME_SERIES_DAILY")
	q.Set("symbol", ticker)
	q.Set("outputsize", "full")
	q.Set("datatype", "csv")
	q.Set("apikey", apiKey)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("requesting data for %s: %w", ticker, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d fetching %s", resp.StatusCode, ticker)
	}

	// TODO: parse resp.Body as CSV into []model.PricePoint
	// Columns from Alpha Vantage: timestamp,open,high,low,close,volume
	// Note: rows come back newest-first — we'll want to decide where
	// reversal happens (here, or downstream in the caller).

	return Ok, resp.Body
}
