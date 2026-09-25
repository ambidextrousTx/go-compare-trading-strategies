// cmd/go-compare-trading-stategies/main.go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"go-compare-trading-strategies/internal/fetch"
	"go-compare-trading-strategies/internal/model"
	"go-compare-trading-strategies/internal/algorithms"
)

func printResult(r model.StrategyResult) {
    fmt.Printf("%s: $%.2f, gain %.2f%% \n", r.Strategy, r.AbsoluteGain, r.PercentReturn)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("loading .env file: %v", err)
	}

	apiKey := os.Getenv("ALPHAVANTAGE_API_KEY")
	if apiKey == "" {
		log.Fatal("ALPHAVANTAGE_API_KEY not set (check your .env file)")
	}

	if len(os.Args) < 2 {
		log.Fatalf("can't run without ticker. Please provide ticker symbol as program argument.")
	}

	ticker := os.Args[1]

	fmt.Println("Ticker is", ticker)

	prices, err := fetch.FetchDaily(ticker, apiKey)
	if err != nil {
		log.Fatalf("fetching prices for %s: %v", ticker, err)
	}

	fmt.Printf("Fetched %d price points for %s\n\n", len(prices), ticker)
	result, err := algorithms.CalculateBuyAndHoldProfit(prices)
	if err != nil {
		log.Fatalf("calculating buy and hold profit for %s: %v", ticker, err)
	}

	printResult(result)

	result2, err := algorithms.CalculateBuyLowSellHighProfit(prices)
	if err != nil {
		log.Fatalf("calculating buy-low-sell-high for %s: %v", ticker, err)
	}

	printResult(result2)
}
