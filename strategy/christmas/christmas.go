package christmas

// Implements the PrizePoolStrategy interface from parent package
// Active from December 20th through January 5th
type Strategy struct{}

// New creates a new instance of the Christmas Bonus strategy
func New() *Strategy {
	return &Strategy{}
}

// Applies a 120% bonus to the base prize pool
// Multiplier: 2.2x (+120% bonus)
func (s *Strategy) CalculatePrizePool(basePrizePool float64) float64 {
	return basePrizePool * 2.2
}

func (s *Strategy) GetStrategyName() string {
	return "Christmas Bonus (120%)"
}
