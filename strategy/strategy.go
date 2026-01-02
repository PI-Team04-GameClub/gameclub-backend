package strategy

type PrizePoolStrategy interface {
	CalculatePrizePool(basePrizePool float64) float64

	GetStrategyName() string
}
