package dtos

import "time"

type CreateTournamentRequest struct {
	Name      string    `json:"name" validate:"required"`
	GameId    uint      `json:"gameId" validate:"required"`
	PrizePool float64   `json:"prizePool" binding:"required,min=0"`
	StartDate time.Time `json:"startDate" binding:"required"`
}

type TournamentResponse struct {
	ID                  uint      `json:"id"`
	Name                string    `json:"name"`
	Game                string    `json:"game"`
	BasePrizePool       float64   `json:"basePrizePool"`
	CalculatedPrizePool float64   `json:"calculatedPrizePool"`
	PrizePool           float64   `json:"prizePool"`
	BonusType           string    `json:"bonusType"`
	BonusMultiplier     float64   `json:"bonusMultiplier"`
	StartDate           time.Time `json:"startDate"`
	Status              string    `json:"status"`
}
