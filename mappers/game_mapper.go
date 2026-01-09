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
		MinPlayers:      game.MinPlayers,
		MaxPlayers:      game.MaxPlayers,
		PlaytimeMinutes: game.PlaytimeMinutes,
		MinAge:          game.MinAge,
		Complexity:      string(game.Complexity),
		Category:        string(game.Category),
		Publisher:       game.Publisher,
		YearPublished:   game.YearPublished,
		Rating:          game.Rating,
	}
}

func ToGameModel(req dtos.CreateGameRequest) models.Game {
	builder := models.NewGameBuilder().
		SetName(req.Name).
		SetDescription(req.Description)

	if req.MinPlayers > 0 && req.MaxPlayers > 0 {
		builder.SetPlayerRange(req.MinPlayers, req.MaxPlayers)
	} else if req.NumberOfPlayers > 0 {
		builder.SetNumberOfPlayers(req.NumberOfPlayers).
			SetMinPlayers(2).
			SetMaxPlayers(req.NumberOfPlayers)
	}

	if req.PlaytimeMinutes > 0 {
		builder.SetPlaytimeMinutes(req.PlaytimeMinutes)
	}

	if req.MinAge > 0 {
		builder.SetMinAge(req.MinAge)
	}

	if req.Complexity != "" {
		builder.SetComplexity(models.GameComplexity(req.Complexity))
	}

	if req.Category != "" {
		builder.SetCategory(models.GameCategory(req.Category))
	}

	if req.Publisher != "" {
		builder.SetPublisher(req.Publisher)
	}

	if req.YearPublished > 0 {
		builder.SetYearPublished(req.YearPublished)
	}

	if req.Rating > 0 {
		builder.SetRating(req.Rating)
	}

	return *builder.Build()
}

func ToGameResponseList(games []models.Game) []dtos.GameResponse {
	responses := make([]dtos.GameResponse, len(games))
	for i, game := range games {
		responses[i] = ToGameResponse(&game)
	}
	return responses
}

func UpdateGameFromRequest(existingGame *models.Game, req dtos.CreateGameRequest) *models.Game {
	builder := models.NewGameBuilder().
		SetID(existingGame.ID).
		SetModel(existingGame.Model).
		SetName(req.Name).
		SetDescription(req.Description)

	if req.MinPlayers > 0 && req.MaxPlayers > 0 {
		builder.SetPlayerRange(req.MinPlayers, req.MaxPlayers)
	} else if req.NumberOfPlayers > 0 {
		builder.SetNumberOfPlayers(req.NumberOfPlayers).
			SetMinPlayers(2).
			SetMaxPlayers(req.NumberOfPlayers)
	}

	if req.PlaytimeMinutes > 0 {
		builder.SetPlaytimeMinutes(req.PlaytimeMinutes)
	}

	if req.MinAge > 0 {
		builder.SetMinAge(req.MinAge)
	}

	if req.Complexity != "" {
		builder.SetComplexity(models.GameComplexity(req.Complexity))
	}

	if req.Category != "" {
		builder.SetCategory(models.GameCategory(req.Category))
	}

	if req.Publisher != "" {
		builder.SetPublisher(req.Publisher)
	}

	if req.YearPublished > 0 {
		builder.SetYearPublished(req.YearPublished)
	}

	if req.Rating > 0 {
		builder.SetRating(req.Rating)
	}

	return builder.Build()
}
