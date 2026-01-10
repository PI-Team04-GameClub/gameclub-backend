package summer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// Given: Nothing (creating a new strategy)

	// When: Creating a new summer strategy
	strategy := New()

	// Then: The strategy should not be nil
	assert.NotNil(t, strategy)
}

func TestStrategy_CalculatePrizePool(t *testing.T) {
	// Given: A summer strategy and a base prize pool of 1000
	strategy := New()
	basePrizePool := 1000.0

	// When: Calculating the prize pool
	result := strategy.CalculatePrizePool(basePrizePool)

	// Then: The result should be 1.2x the base prize pool (20% bonus)
	assert.Equal(t, 1200.0, result)
}

func TestStrategy_CalculatePrizePool_Zero(t *testing.T) {
	// Given: A summer strategy and a base prize pool of 0
	strategy := New()
	basePrizePool := 0.0

	// When: Calculating the prize pool
	result := strategy.CalculatePrizePool(basePrizePool)

	// Then: The result should be zero
	assert.Equal(t, 0.0, result)
}

func TestStrategy_CalculatePrizePool_Large(t *testing.T) {
	// Given: A summer strategy and a large base prize pool
	strategy := New()
	basePrizePool := 100000.0

	// When: Calculating the prize pool
	result := strategy.CalculatePrizePool(basePrizePool)

	// Then: The result should be 1.2x the base prize pool
	assert.Equal(t, 120000.0, result)
}

func TestStrategy_GetStrategyName(t *testing.T) {
	// Given: A summer strategy
	strategy := New()

	// When: Getting the strategy name
	name := strategy.GetStrategyName()

	// Then: The name should be "Summer Bonus (20%)"
	assert.Equal(t, "Summer Bonus (20%)", name)
}

func TestStrategy_CalculatePrizePool_Decimal(t *testing.T) {
	// Given: A summer strategy and a base prize pool with decimals
	strategy := New()
	basePrizePool := 1500.50

	// When: Calculating the prize pool
	result := strategy.CalculatePrizePool(basePrizePool)

	// Then: The result should be 1.2x the base prize pool
	assert.Equal(t, 1800.60, result)
}
