package strategy

import (
	"time"

	"github.com/PI-Team04-GameClub/gameclub-backend/strategy/christmas"
	"github.com/PI-Team04-GameClub/gameclub-backend/strategy/normal"
	"github.com/PI-Team04-GameClub/gameclub-backend/strategy/summer"
)

// Applies the prize pool calculation strategy
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

// GetCurrentStrategy returns the currently active strategy
func (c *Calculator) GetCurrentStrategy() PrizePoolStrategy {
	return c.strategy
}

// GetStrategyForDate returns the appropriate strategy based on the given date
// This is a factory function that selects the correct strategy implementation
//
// Rules:
// - July 1-31: Summer Bonus Strategy (1.2x multiplier)
// - December 20 - January 5: Christmas Bonus Strategy (2.2x multiplier)
// - All other dates: Normal Strategy (1.0x multiplier)
func GetStrategyForDate(date time.Time) PrizePoolStrategy {
	month := date.Month()
	day := date.Day()

	// Summer Bonus: July (entire month)
	if month == time.July {
		return summer.New()
	}

	// Christmas Bonus: December 20th to January 5th
	if (month == time.December && day >= 20) || (month == time.January && day <= 5) {
		return christmas.New()
	}

	// Normal: All other times
	return normal.New()
}

// GetStrategyForNow returns the appropriate strategy based on current date/time
func GetStrategyForNow() PrizePoolStrategy {
	return GetStrategyForDate(time.Now())
}
