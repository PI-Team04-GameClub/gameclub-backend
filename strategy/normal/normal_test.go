package normal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// Given: Nothing (creating a new strategy)

	// When: Creating a new normal strategy
	strategy := New()

	// Then: The strategy should not be nil
	assert.NotNil(t, strategy)
}

func TestStrategy_CalculatePrizePool(t *testing.T) {
	// Given: A normal strategy and a base prize pool of 1000
	strategy := New()
	basePrizePool := 1000.0

	// When: Calculating the prize pool
	result := strategy.CalculatePrizePool(basePrizePool)

	// Then: The result should be the same as the base prize pool (no bonus)
	assert.Equal(t, 1000.0, result)
}

func TestStrategy_CalculatePrizePool_Zero(t *testing.T) {
	// Given: A normal strategy and a base prize pool of 0
	strategy := New()
	basePrizePool := 0.0

	// When: Calculating the prize pool
	result := strategy.CalculatePrizePool(basePrizePool)

	// Then: The result should be zero
	assert.Equal(t, 0.0, result)
}

func TestStrategy_CalculatePrizePool_Large(t *testing.T) {
	// Given: A normal strategy and a large base prize pool
	strategy := New()
	basePrizePool := 999999.99

	// When: Calculating the prize pool
	result := strategy.CalculatePrizePool(basePrizePool)

	// Then: The result should be the same as the base prize pool
	assert.Equal(t, 999999.99, result)
}

func TestStrategy_GetStrategyName(t *testing.T) {
	// Given: A normal strategy
	strategy := New()

	// When: Getting the strategy name
	name := strategy.GetStrategyName()

	// Then: The name should be "Normal"
	assert.Equal(t, "Normal", name)
}
