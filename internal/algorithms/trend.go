package algorithms

import (
	"go-compare-trading-strategies/internal/model"
	"fmt"
)

type direction int

func (d direction) String() string {
	switch d {
	case up: return "Up"
	case down: return "Down"
	default: return "Unknown"
	}
}

const (
	unknown direction = iota
	up
	down
)

type maxima int
const (
	neutral maxima = iota
	localMax
	localMin
)

func (m maxima) String() string {
	switch m {
	case localMax: return "LocalMax"
	case localMin: return "LocalMin"
	default: return "Unknown"
	}
}

type TurningPoint struct {
	Point model.PricePoint
	Type maxima
}

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

func FindTurningPoints(prices []model.PricePoint, trend []direction) ([]TurningPoint, error) {
	turningPoints := make([]TurningPoint, 0, 500)

	for i := 1; i < len(trend); i++ {
		switch {
		case trend[i] == up && trend[i - 1] == down:
			turningPoints = append(turningPoints, TurningPoint{Point: prices[i], Type: localMin})
		case trend[i] == down && trend[i - 1] == up:
			turningPoints = append(turningPoints, TurningPoint{Point: prices[i], Type: localMax})
		}
	}

	return turningPoints, nil

}

func CalculateBuyLowSellHighProfit(prices []model.PricePoint) (model.StrategyResult, error) {

	trend, err := CalculateTrend(prices)
	if err != nil {
		return model.StrategyResult{}, fmt.Errorf("calculating trend: %w", err)
	}

	turningPoints, err := FindTurningPoints(prices, trend)
	if err != nil {
		return model.StrategyResult{}, fmt.Errorf("calculating turning points: %w", err)
	}

	first := prices[0].Close
	gain := 0.0
	holding := false
	buyPrice := 0.0

	for _, tp := range turningPoints {
		switch tp.Type {
		case localMin:
			if !holding {
				buyPrice = tp.Point.Close
				holding = true
			}
		case localMax:
			if holding {
				gain += tp.Point.Close - buyPrice
				holding = false
			}
		}
	}

	return model.StrategyResult{
		Strategy:      "Buy Local Min Sell Local Max",
		AbsoluteGain:  gain,
		PercentReturn: (gain / first) * 100,
	}, nil

}
