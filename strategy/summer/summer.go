package summer

// Implements the PrizePoolStrategy interface from parent package
// Active during the entire month of July
type Strategy struct{}

// New creates a new instance of the Summer Bonus strategy
func New() *Strategy {
	return &Strategy{}
}

// Applies a 20% bonus to the base prize pool
// Multiplier: 1.2x (+20% bonus)
func (s *Strategy) CalculatePrizePool(basePrizePool float64) float64 {
	return basePrizePool * 1.2
}

// GetStrategyName returns the display name of this strategy
func (s *Strategy) GetStrategyName() string {
	return "Summer Bonus (20%)"
}
