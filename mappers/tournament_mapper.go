package mappers

import (
	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
)

func ToTournamentResponse(tournament *models.Tournament) dtos.TournamentResponse {
	return dtos.TournamentResponse{
		ID:        tournament.ID,
		Name:      tournament.Name,
		Game:      tournament.Game.Name,
		PrizePool: tournament.PrizePool,
		StartDate: tournament.StartDate,
		Status:    string(tournament.Status),
	}
}

func ToTournamentModel(req dtos.CreateTournamentRequest) models.Tournament {
	return models.Tournament{
		Name:      req.Name,
		PrizePool: req.PrizePool,
		GameID:    req.GameId,
		StartDate: req.StartDate,
		Status:    models.StatusUpcoming,
	}
}

func ToTournamentResponseList(tournaments []models.Tournament) []dtos.TournamentResponse {
	responses := make([]dtos.TournamentResponse, len(tournaments))
	for i, tournament := range tournaments {
		responses[i] = ToTournamentResponse(&tournament)
	}
	return responses
}
