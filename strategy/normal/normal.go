package normal

// Strategy represents the normal prize pool calculation (no bonus)
// Implements the PrizePoolStrategy interface from parent package
type Strategy struct{}

// New creates a new instance of the Normal strategy
func New() *Strategy {
	return &Strategy{}
}

// Returns the base prize pool without any modifications
// Multiplier: 1.0x (no bonus)
func (s *Strategy) CalculatePrizePool(basePrizePool float64) float64 {
	return basePrizePool
}

// GetStrategyName returns the display name of this strategy
func (s *Strategy) GetStrategyName() string {
	return "Normal"
}
