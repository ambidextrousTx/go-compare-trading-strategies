// internal/fetch/alphavantage.go
package fetch

import (
	"fmt"
	"net/http"
	"net/url"
	"bufio"
	"log"
	"strings"
	"time"
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

	fmt.Println("Requesting", u)
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("requesting data for %s: %w", ticker, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d fetching %s", resp.StatusCode, ticker)
	}

	// fmt.Printf("%s", body)
	// TODO: parse resp.Body as CSV into []model.PricePoint
	// Columns from Alpha Vantage: timestamp,open,high,low,close,volume
	// Note: rows come back newest-first — we'll want to decide where
	// reversal happens (here, or downstream in the caller).

	pricePoints := make([]model.PricePoint, 0, 500)

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		lineElements := strings.Split(scanner.Text(), ",")
		fmt.Println("Processing", lineElements)
		date, _ := time.Parse(time.RFC3339, lineElements[0])
		open, _ := strconv.ParseFloat(lineElements[1], 64)
		high, _ := strconv.ParseFloat(lineElements[2], 64)
		low, _ := strconv.ParseFloat(lineElements[3], 64)
		closing, _ := strconv.ParseFloat(lineElements[4], 64)
		volume, _ := strconv.ParseInt(lineElements[5], 10, 64)

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

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return pricePoints, nil
}
