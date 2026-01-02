package mappers

import (
	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
)

func ToTournamentResponse(tournament *models.Tournament) dtos.TournamentResponse {
	bonusMultiplier := 1.0
	if tournament.BasePrizePool > 0 {
		bonusMultiplier = tournament.CalculatedPrizePool / tournament.BasePrizePool
	}

	return dtos.TournamentResponse{
		ID:                  tournament.ID,
		Name:                tournament.Name,
		Game:                tournament.Game.Name,
		BasePrizePool:       tournament.BasePrizePool,
		CalculatedPrizePool: tournament.CalculatedPrizePool,
		PrizePool:           tournament.CalculatedPrizePool,
		BonusType:           tournament.BonusType,
		BonusMultiplier:     bonusMultiplier,
		StartDate:           tournament.StartDate,
		Status:              string(tournament.Status),
	}
}

func ToTournamentModel(req dtos.CreateTournamentRequest) models.Tournament {
	tournament := models.Tournament{
		Name:          req.Name,
		GameID:        req.GameId,
		BasePrizePool: req.PrizePool,
		StartDate:     req.StartDate,
		Status:        models.StatusUpcoming,
	}

	// Apply the Strategy pattern to calculate prize pool based on start date
	tournament.ApplyPrizePoolStrategy()

	return tournament
}

func ToTournamentResponseList(tournaments []models.Tournament) []dtos.TournamentResponse {
	responses := make([]dtos.TournamentResponse, len(tournaments))
	for i, tournament := range tournaments {
		responses[i] = ToTournamentResponse(&tournament)
	}
	return responses
}

// UpdateTournamentFromRequest updates an existing tournament and recalculates prize pool if date changed
func UpdateTournamentFromRequest(existingTournament *models.Tournament, req dtos.CreateTournamentRequest) *models.Tournament {
	existingTournament.Name = req.Name
	existingTournament.GameID = req.GameId
	existingTournament.BasePrizePool = req.PrizePool

	// Check if start date changed
	dateChanged := !existingTournament.StartDate.Equal(req.StartDate)
	existingTournament.StartDate = req.StartDate

	// Recalculate prize pool if base amount or date changed
	if dateChanged || existingTournament.BasePrizePool != req.PrizePool {
		existingTournament.ApplyPrizePoolStrategy()
	}

	return existingTournament
}
