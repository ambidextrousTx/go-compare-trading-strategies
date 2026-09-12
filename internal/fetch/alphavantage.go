// internal/fetch/alphavantage.go
package fetch

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"time"

	"go-compare-trading-strategies/internal/model"
)

const baseURL = "https://www.alphavantage.co/query"

func parseIntColumn(record []string, index int, name string) (int64, error) {
	v, err := strconv.ParseInt(record[index], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing %s (%q): %w", name, record[index], err)
	}
	return v, nil
}

func parseFloatColumn(record []string, index int, name string) (float64, error) {
	v, err := strconv.ParseFloat(record[index], 64)
	if err != nil {
		return 0, fmt.Errorf("parsing %s (%q): %w", name, record[index], err)
	}
	return v, nil
}

// FetchDaily retrieves full daily price history for a ticker.
func FetchDaily(ticker, apiKey string) ([]model.PricePoint, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("parsing base URL: %w", err)
	}

	q := u.Query()
	q.Set("function", "TIME_SERIES_DAILY")
	q.Set("symbol", ticker)
	// q.Set("outputsize", "full") // Premium feature
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

	// Columns from Alpha Vantage: timestamp,open,high,low,close,volume
	// Note: rows come back newest-first — we'll want to decide where
	// reversal happens (here, or downstream in the caller).

	pricePoints := make([]model.PricePoint, 0, 500)

	reader := csv.NewReader(resp.Body)

	header, err := reader.Read() // consume the header
	if err != nil {
		return nil, fmt.Errorf("reading CSV header: %w", err)
	}

	_ = header

	for {
		record, err := reader.Read()

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("reading CSV row: %w", err)
		}

		date, err := time.Parse(time.DateOnly, record[0])
		if err != nil {
			return nil, err
		}
		open, err := parseFloatColumn(record, 1, "open")
		if err != nil {
			return nil, fmt.Errorf("parsing date (%q): %w", record[0], err)
		}

		high, err := parseFloatColumn(record, 2, "high")
		if err != nil {
			return nil, err
		}

		low, err := parseFloatColumn(record, 3, "low")
		if err != nil {
			return nil, err
		}

		closing, err := parseFloatColumn(record, 4, "closing")
		if err != nil {
			return nil, err
		}

		volume, err := parseIntColumn(record, 5, "volume")
		if err != nil {
			return nil, err
		}

		pricePoints = append(pricePoints,
			model.PricePoint{
				Date:   date,
				Open:   open,
				High:   high,
				Low:    low,
				Close:  closing,
				Volume: volume,
			})
	}

	slices.Reverse(pricePoints)
	return pricePoints, nil
}
