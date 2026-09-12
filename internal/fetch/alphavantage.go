// internal/fetch/alphavantage.go
package fetch

import (
	"io"
	"encoding/csv"
	"fmt"
	"time"
	"net/http"
	"net/url"
	"strconv"

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
		return nil, fmt.Errorf("Reading CSV header: %w", header)
	}

	_ = header

	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
        return nil, fmt.Errorf("Reading CSV row: %w", err)
    }

		date, err := time.Parse(time.UnixDate, record[0])
		if err != nil {
			fmt.Errorf("Error reading row: %w", record)
		}
		open, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			fmt.Errorf("Error reading row: %w", record)
		}
		high, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			fmt.Errorf("Error reading row: %w", record)
		}
		low, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			fmt.Errorf("Error reading row: %w", record)
		}
		closing, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			fmt.Errorf("Error reading row: %w", record)
		}
		volume, err := strconv.ParseInt(record[5], 10, 64)
		if err != nil {
			fmt.Errorf("Error reading row: %w", record)
		}

		pricePoints = append(pricePoints,
			model.PricePoint {
				Date:   date,
				Open:   open,
				High:   high,
				Low:    low,
				Close:  closing,
				Volume: volume,
			})
	}

	return pricePoints, nil
}
