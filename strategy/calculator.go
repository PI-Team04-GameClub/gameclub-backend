package strategy

import (
	"time"

	"github.com/PI-Team04-GameClub/gameclub-backend/strategy/christmas"
	"github.com/PI-Team04-GameClub/gameclub-backend/strategy/normal"
	"github.com/PI-Team04-GameClub/gameclub-backend/strategy/summer"
)

type Calculator struct {
	strategy PrizePoolStrategy
}

func NewCalculator(strategy PrizePoolStrategy) *Calculator {
	return &Calculator{
		strategy: strategy,
	}
}

func (c *Calculator) SetStrategy(strategy PrizePoolStrategy) {
	c.strategy = strategy
}

func (c *Calculator) Calculate(basePrizePool float64) float64 {
	return c.strategy.CalculatePrizePool(basePrizePool)
}

func (c *Calculator) GetCurrentStrategy() PrizePoolStrategy {
	return c.strategy
}

func GetStrategyForDate(date time.Time) PrizePoolStrategy {
	month := date.Month()
	day := date.Day()

	if month == time.July {
		return summer.New()
	}

	if (month == time.December && day >= 20) || (month == time.January && day <= 5) {
		return christmas.New()
	}

	return normal.New()
}

func GetStrategyForNow() PrizePoolStrategy {
	return GetStrategyForDate(time.Now())
}
