package christmas

type Strategy struct{}

func New() *Strategy {
	return &Strategy{}
}

func (s *Strategy) CalculatePrizePool(basePrizePool float64) float64 {
	return basePrizePool * 2.2
}

func (s *Strategy) GetStrategyName() string {
	return "Christmas Bonus (120%)"
}
