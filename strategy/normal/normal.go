package normal

type Strategy struct{}

func New() *Strategy {
	return &Strategy{}
}

func (s *Strategy) CalculatePrizePool(basePrizePool float64) float64 {
	return basePrizePool
}

func (s *Strategy) GetStrategyName() string {
	return "Normal"
}
