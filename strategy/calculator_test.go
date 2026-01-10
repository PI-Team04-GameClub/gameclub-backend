package strategy

import (
	"testing"
	"time"

	"github.com/PI-Team04-GameClub/gameclub-backend/strategy/christmas"
	"github.com/PI-Team04-GameClub/gameclub-backend/strategy/normal"
	"github.com/PI-Team04-GameClub/gameclub-backend/strategy/summer"
	"github.com/stretchr/testify/assert"
)

func TestNewCalculator(t *testing.T) {
	// Given: A normal strategy
	strategy := normal.New()

	// When: Creating a new calculator with the strategy
	calculator := NewCalculator(strategy)

	// Then: The calculator should be created with the correct strategy
	assert.NotNil(t, calculator)
	assert.Equal(t, strategy, calculator.GetCurrentStrategy())
}

func TestCalculator_Calculate_NormalStrategy(t *testing.T) {
	// Given: A calculator with normal strategy and a base prize pool of 1000
	strategy := normal.New()
	calculator := NewCalculator(strategy)
	basePrizePool := 1000.0

	// When: Calculating the prize pool
	result := calculator.Calculate(basePrizePool)

	// Then: The result should be the same as the base prize pool (1.0x multiplier)
	assert.Equal(t, 1000.0, result)
}

func TestCalculator_Calculate_SummerStrategy(t *testing.T) {
	// Given: A calculator with summer strategy and a base prize pool of 1000
	strategy := summer.New()
	calculator := NewCalculator(strategy)
	basePrizePool := 1000.0

	// When: Calculating the prize pool
	result := calculator.Calculate(basePrizePool)

	// Then: The result should be 1.2x the base prize pool (20% bonus)
	assert.Equal(t, 1200.0, result)
}

func TestCalculator_Calculate_ChristmasStrategy(t *testing.T) {
	// Given: A calculator with Christmas strategy and a base prize pool of 1000
	strategy := christmas.New()
	calculator := NewCalculator(strategy)
	basePrizePool := 1000.0

	// When: Calculating the prize pool
	result := calculator.Calculate(basePrizePool)

	// Then: The result should be 2.2x the base prize pool (120% bonus)
	assert.Equal(t, 2200.0, result)
}

func TestCalculator_SetStrategy(t *testing.T) {
	// Given: A calculator with normal strategy
	normalStrategy := normal.New()
	calculator := NewCalculator(normalStrategy)

	// When: Setting a new summer strategy
	summerStrategy := summer.New()
	calculator.SetStrategy(summerStrategy)

	// Then: The calculator should use the new strategy
	assert.Equal(t, summerStrategy, calculator.GetCurrentStrategy())
	assert.Equal(t, 1200.0, calculator.Calculate(1000.0))
}

func TestCalculator_Calculate_ZeroPrizePool(t *testing.T) {
	// Given: A calculator with summer strategy and zero prize pool
	strategy := summer.New()
	calculator := NewCalculator(strategy)

	// When: Calculating with zero prize pool
	result := calculator.Calculate(0.0)

	// Then: The result should be zero
	assert.Equal(t, 0.0, result)
}

func TestCalculator_Calculate_LargePrizePool(t *testing.T) {
	// Given: A calculator with Christmas strategy and a large prize pool
	strategy := christmas.New()
	calculator := NewCalculator(strategy)
	largePrizePool := 1000000.0

	// When: Calculating the prize pool
	result := calculator.Calculate(largePrizePool)

	// Then: The result should be correctly calculated
	assert.Equal(t, 2200000.0, result)
}

func TestGetStrategyForDate_July(t *testing.T) {
	// Given: A date in July (summer month)
	julyDate := time.Date(2024, time.July, 15, 0, 0, 0, 0, time.UTC)

	// When: Getting the strategy for that date
	strategy := GetStrategyForDate(julyDate)

	// Then: The summer strategy should be returned
	assert.Equal(t, "Summer Bonus (20%)", strategy.GetStrategyName())
}

func TestGetStrategyForDate_December25(t *testing.T) {
	// Given: A date on Christmas day (December 25)
	christmasDate := time.Date(2024, time.December, 25, 0, 0, 0, 0, time.UTC)

	// When: Getting the strategy for that date
	strategy := GetStrategyForDate(christmasDate)

	// Then: The Christmas strategy should be returned
	assert.Equal(t, "Christmas Bonus (120%)", strategy.GetStrategyName())
}

func TestGetStrategyForDate_January3(t *testing.T) {
	// Given: A date in early January (Christmas period continues)
	newYearDate := time.Date(2024, time.January, 3, 0, 0, 0, 0, time.UTC)

	// When: Getting the strategy for that date
	strategy := GetStrategyForDate(newYearDate)

	// Then: The Christmas strategy should be returned
	assert.Equal(t, "Christmas Bonus (120%)", strategy.GetStrategyName())
}

func TestGetStrategyForDate_February(t *testing.T) {
	// Given: A date in February (normal period)
	date := time.Date(2024, time.February, 15, 0, 0, 0, 0, time.UTC)

	// When: Getting the strategy for that date
	strategy := GetStrategyForDate(date)

	// Then: The normal strategy should be returned
	assert.Equal(t, "Normal", strategy.GetStrategyName())
}

func TestGetStrategyForDate_March(t *testing.T) {
	// Given: A date in March (normal period)
	date := time.Date(2024, time.March, 15, 0, 0, 0, 0, time.UTC)

	// When: Getting the strategy for that date
	strategy := GetStrategyForDate(date)

	// Then: The normal strategy should be returned
	assert.Equal(t, "Normal", strategy.GetStrategyName())
}

func TestGetStrategyForDate_April(t *testing.T) {
	// Given: A date in April (normal period)
	date := time.Date(2024, time.April, 15, 0, 0, 0, 0, time.UTC)

	// When: Getting the strategy for that date
	strategy := GetStrategyForDate(date)

	// Then: The normal strategy should be returned
	assert.Equal(t, "Normal", strategy.GetStrategyName())
}

func TestGetStrategyForDate_May(t *testing.T) {
	// Given: A date in May (normal period)
	date := time.Date(2024, time.May, 15, 0, 0, 0, 0, time.UTC)

	// When: Getting the strategy for that date
	strategy := GetStrategyForDate(date)

	// Then: The normal strategy should be returned
	assert.Equal(t, "Normal", strategy.GetStrategyName())
}

func TestGetStrategyForDate_June(t *testing.T) {
	// Given: A date in June (normal period)
	date := time.Date(2024, time.June, 15, 0, 0, 0, 0, time.UTC)

	// When: Getting the strategy for that date
	strategy := GetStrategyForDate(date)

	// Then: The normal strategy should be returned
	assert.Equal(t, "Normal", strategy.GetStrategyName())
}

func TestGetStrategyForDate_August(t *testing.T) {
	// Given: A date in August (normal period)
	date := time.Date(2024, time.August, 15, 0, 0, 0, 0, time.UTC)

	// When: Getting the strategy for that date
	strategy := GetStrategyForDate(date)

	// Then: The normal strategy should be returned
	assert.Equal(t, "Normal", strategy.GetStrategyName())
}

func TestGetStrategyForDate_September(t *testing.T) {
	// Given: A date in September (normal period)
	date := time.Date(2024, time.September, 15, 0, 0, 0, 0, time.UTC)

	// When: Getting the strategy for that date
	strategy := GetStrategyForDate(date)

	// Then: The normal strategy should be returned
	assert.Equal(t, "Normal", strategy.GetStrategyName())
}

func TestGetStrategyForDate_October(t *testing.T) {
	// Given: A date in October (normal period)
	date := time.Date(2024, time.October, 15, 0, 0, 0, 0, time.UTC)

	// When: Getting the strategy for that date
	strategy := GetStrategyForDate(date)

	// Then: The normal strategy should be returned
	assert.Equal(t, "Normal", strategy.GetStrategyName())
}

func TestGetStrategyForDate_November(t *testing.T) {
	// Given: A date in November (normal period)
	date := time.Date(2024, time.November, 15, 0, 0, 0, 0, time.UTC)

	// When: Getting the strategy for that date
	strategy := GetStrategyForDate(date)

	// Then: The normal strategy should be returned
	assert.Equal(t, "Normal", strategy.GetStrategyName())
}

func TestGetStrategyForDate_EarlyDecember(t *testing.T) {
	// Given: A date on December 10 (before Christmas period)
	date := time.Date(2024, time.December, 10, 0, 0, 0, 0, time.UTC)

	// When: Getting the strategy for that date
	strategy := GetStrategyForDate(date)

	// Then: The normal strategy should be returned
	assert.Equal(t, "Normal", strategy.GetStrategyName())
}

func TestGetStrategyForDate_LateJanuary(t *testing.T) {
	// Given: A date on January 10 (after Christmas period)
	date := time.Date(2024, time.January, 10, 0, 0, 0, 0, time.UTC)

	// When: Getting the strategy for that date
	strategy := GetStrategyForDate(date)

	// Then: The normal strategy should be returned
	assert.Equal(t, "Normal", strategy.GetStrategyName())
}

func TestGetStrategyForDate_December19(t *testing.T) {
	// Given: December 19th (just before Christmas period)
	dec19 := time.Date(2024, time.December, 19, 0, 0, 0, 0, time.UTC)

	// When: Getting the strategy
	strategy := GetStrategyForDate(dec19)

	// Then: Normal strategy should be returned
	assert.Equal(t, "Normal", strategy.GetStrategyName())
}

func TestGetStrategyForDate_December20(t *testing.T) {
	// Given: December 20th (start of Christmas period)
	dec20 := time.Date(2024, time.December, 20, 0, 0, 0, 0, time.UTC)

	// When: Getting the strategy
	strategy := GetStrategyForDate(dec20)

	// Then: Christmas strategy should be returned
	assert.Equal(t, "Christmas Bonus (120%)", strategy.GetStrategyName())
}

func TestGetStrategyForDate_January5(t *testing.T) {
	// Given: January 5th (last day of Christmas period)
	jan5 := time.Date(2024, time.January, 5, 0, 0, 0, 0, time.UTC)

	// When: Getting the strategy
	strategy := GetStrategyForDate(jan5)

	// Then: Christmas strategy should be returned
	assert.Equal(t, "Christmas Bonus (120%)", strategy.GetStrategyName())
}

func TestGetStrategyForDate_January6(t *testing.T) {
	// Given: January 6th (day after Christmas period)
	jan6 := time.Date(2024, time.January, 6, 0, 0, 0, 0, time.UTC)

	// When: Getting the strategy
	strategy := GetStrategyForDate(jan6)

	// Then: Normal strategy should be returned
	assert.Equal(t, "Normal", strategy.GetStrategyName())
}

func TestGetStrategyForNow(t *testing.T) {
	// Given: The current time
	now := time.Now()

	// When: Getting the strategy for now
	strategy := GetStrategyForNow()

	// Then: The strategy should match what GetStrategyForDate returns for the current time
	expectedStrategy := GetStrategyForDate(now)
	assert.Equal(t, expectedStrategy.GetStrategyName(), strategy.GetStrategyName())
}
