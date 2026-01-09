package summer

type Strategy struct{}

func New() *Strategy {
	return &Strategy{}
}

func (s *Strategy) CalculatePrizePool(basePrizePool float64) float64 {
	return basePrizePool * 1.2
}

func (s *Strategy) GetStrategyName() string {
	return "Summer Bonus (20%)"
}
