package mappers

import (
	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
)

func ToTeamResponse(team *models.Team) dtos.TeamResponse {
	return dtos.TeamResponse{
		ID:   team.ID,
		Name: team.Name,
	}
}

func ToTeamModel(req dtos.CreateTeamRequest) models.Team {
	return models.Team{
		Name: req.Name,
	}
}

func ToTeamResponseList(teams []models.Team) []dtos.TeamResponse {
	responses := make([]dtos.TeamResponse, len(teams))
	for i, team := range teams {
		responses[i] = ToTeamResponse(&team)
	}
	return responses
}
