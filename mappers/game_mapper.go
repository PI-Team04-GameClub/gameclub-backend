package mappers

import (
	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
)

func ToGameResponse(game *models.Game) dtos.GameResponse {
	return dtos.GameResponse{
		ID:              game.ID,
		Name:            game.Name,
		Description:     game.Description,
		NumberOfPlayers: game.NumberOfPlayers,
	}
}

func ToGameModel(req dtos.CreateGameRequest) models.Game {
	return models.Game{
		Name:            req.Name,
		Description:     req.Description,
		NumberOfPlayers: req.NumberOfPlayers,
	}
}

func ToGameResponseList(games []models.Game) []dtos.GameResponse {
	responses := make([]dtos.GameResponse, len(games))
	for i, game := range games {
		responses[i] = ToGameResponse(&game)
	}
	return responses
}
