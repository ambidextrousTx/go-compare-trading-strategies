package algorithms

import (
	"fmt"
	"go-compare-trading-strategies/internal/model"
)

func CalculateBuyAndHoldProfit(prices []model.PricePoint) (model.StrategyResult, error) {
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

