package algorithms

import (
	"fmt"
	"go-compare-trading-strategies/internal/model"
)

type direction int

const (
    unknown direction = iota
    up
    down
)

func CalculateTrend(prices []model.PricePoint) ([]direction, error) {
	trend := make([]direction, 0, 500)

	for i := 1; i < len(prices); i++ {
		switch {
		case prices[i].Close > prices[i-1].Close:
			trend = append(trend, up)
		case prices[i].Close < prices[i-1].Close:
			trend = append(trend, down)
		default:
			trend = append(trend, unknown)
		}
	}

	return trend, nil
}
