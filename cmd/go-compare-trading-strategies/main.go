// cmd/go-compare-trading-stategies/main.go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"go-compare-trading-strategies/internal/fetch"
	"go-compare-trading-strategies/internal/model"
)

func calculateBuyAndHoldProfit(prices []model.PricePoint) (model.StrategyResult, error) {
	if len(prices) == 0 {
		return model.StrategyResult{}, fmt.Errorf("no price data to evaluate")
	}

	first := prices[0].Close
    last := prices[len(prices)-1].Close
    gain := last - first

	return model.StrategyResult{
		Strategy:      "Buy and Hold",
		AbsoluteGain:  gain,
		PercentReturn: (gain / first) * 100,
	}, nil
}

func printResult(r model.StrategyResult) {
    fmt.Printf("%s: $%.2f, gain %.2f%% \n\n", r.Strategy, r.AbsoluteGain, r.PercentReturn)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("loading .env file: %v", err)
	}

	apiKey := os.Getenv("ALPHAVANTAGE_API_KEY")
	if apiKey == "" {
		log.Fatal("ALPHAVANTAGE_API_KEY not set (check your .env file)")
	}

	// TODO: take ticker as a CLI arg instead of hardcoding
	ticker := "AAPL"
	fmt.Println("Ticker is", ticker)

	prices, err := fetch.FetchDaily(ticker, apiKey)
	if err != nil {
		log.Fatalf("fetching prices for %s: %v", ticker, err)
	}

	fmt.Printf("fetched %d price points for %s\n", len(prices), ticker)
	result, err := calculateBuyAndHoldProfit(prices)
	if err != nil {
		log.Fatalf("calculating buy and hold profit for %s: %v", ticker, err)
	}

	printResult(result)
}
