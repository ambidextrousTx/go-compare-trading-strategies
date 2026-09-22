package algorithms

import (
	"go-compare-trading-strategies/internal/model"
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
