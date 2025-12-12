package dtos

import "time"

type CreateTournamentRequest struct {
	Name      string    `json:"name" validate:"required"`
	GameId    uint      `json:"gameId" validate:"required"`
	PrizePool float64   `json:"prizePool" binding:"required,min=0"`
	StartDate time.Time `json:"startDate" binding:"required"`
}

type TournamentResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Game      string    `json:"game"`
	PrizePool float64   `json:"prizePool" binding:"min=0"`
	StartDate time.Time `json:"startDate"`
	Status    string    `json:"status"`
}
