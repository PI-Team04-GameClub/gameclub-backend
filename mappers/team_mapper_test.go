package mappers

import (
	"testing"

	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestToTeamResponse_BasicFields(t *testing.T) {
	// Given: A team model with all fields populated
	team := &models.Team{
		Model: gorm.Model{ID: 1},
		Name:  "Alpha Team",
	}

	// When: Converting to response
	response := ToTeamResponse(team)

	// Then: All fields should be correctly mapped
	assert.Equal(t, uint(1), response.ID)
	assert.Equal(t, "Alpha Team", response.Name)
}

func TestToTeamResponse_EmptyName(t *testing.T) {
	// Given: A team model with empty name
	team := &models.Team{
		Model: gorm.Model{ID: 2},
		Name:  "",
	}

	// When: Converting to response
	response := ToTeamResponse(team)

	// Then: Empty name should be preserved
	assert.Equal(t, uint(2), response.ID)
	assert.Equal(t, "", response.Name)
}

func TestToTeamResponse_LongName(t *testing.T) {
	// Given: A team model with a long name
	longName := "This Is A Very Long Team Name That Should Still Work Correctly"
	team := &models.Team{
		Model: gorm.Model{ID: 3},
		Name:  longName,
	}

	// When: Converting to response
	response := ToTeamResponse(team)

	// Then: Long name should be preserved
	assert.Equal(t, longName, response.Name)
}

func TestToTeamModel_BasicRequest(t *testing.T) {
	// Given: A create team request
	req := dtos.CreateTeamRequest{
		Name: "Beta Team",
	}

	// When: Converting to model
	team := ToTeamModel(req)

	// Then: The name should be set correctly
	assert.Equal(t, "Beta Team", team.Name)
}

func TestToTeamModel_EmptyName(t *testing.T) {
	// Given: A create team request with empty name
	req := dtos.CreateTeamRequest{
		Name: "",
	}

	// When: Converting to model
	team := ToTeamModel(req)

	// Then: Empty name should be preserved
	assert.Equal(t, "", team.Name)
}

func TestToTeamModel_SpecialCharacters(t *testing.T) {
	// Given: A create team request with special characters
	req := dtos.CreateTeamRequest{
		Name: "Team #1 - The Champions!",
	}

	// When: Converting to model
	team := ToTeamModel(req)

	// Then: Special characters should be preserved
	assert.Equal(t, "Team #1 - The Champions!", team.Name)
}

func TestToTeamResponseList_EmptySlice(t *testing.T) {
	// Given: An empty slice of teams
	teams := []models.Team{}

	// When: Converting to response list
	responses := ToTeamResponseList(teams)

	// Then: The result should be an empty slice
	assert.Empty(t, responses)
	assert.Len(t, responses, 0)
}

func TestToTeamResponseList_SingleTeam(t *testing.T) {
	// Given: A slice with one team
	teams := []models.Team{
		{Model: gorm.Model{ID: 1}, Name: "Solo Team"},
	}

	// When: Converting to response list
	responses := ToTeamResponseList(teams)

	// Then: One team should be converted
	assert.Len(t, responses, 1)
	assert.Equal(t, "Solo Team", responses[0].Name)
}

func TestToTeamResponseList_MultipleTeams(t *testing.T) {
	// Given: A slice with multiple teams
	teams := []models.Team{
		{Model: gorm.Model{ID: 1}, Name: "Team Alpha"},
		{Model: gorm.Model{ID: 2}, Name: "Team Beta"},
		{Model: gorm.Model{ID: 3}, Name: "Team Gamma"},
		{Model: gorm.Model{ID: 4}, Name: "Team Delta"},
	}

	// When: Converting to response list
	responses := ToTeamResponseList(teams)

	// Then: All teams should be converted
	assert.Len(t, responses, 4)
	assert.Equal(t, "Team Alpha", responses[0].Name)
	assert.Equal(t, "Team Beta", responses[1].Name)
	assert.Equal(t, "Team Gamma", responses[2].Name)
	assert.Equal(t, "Team Delta", responses[3].Name)
}

func TestToTeamResponseList_PreservesOrder(t *testing.T) {
	// Given: A slice with teams in specific order
	teams := []models.Team{
		{Model: gorm.Model{ID: 5}, Name: "Zebra"},
		{Model: gorm.Model{ID: 3}, Name: "Apple"},
		{Model: gorm.Model{ID: 1}, Name: "Mango"},
	}

	// When: Converting to response list
	responses := ToTeamResponseList(teams)

	// Then: The order should be preserved (not sorted)
	assert.Equal(t, "Zebra", responses[0].Name)
	assert.Equal(t, "Apple", responses[1].Name)
	assert.Equal(t, "Mango", responses[2].Name)
}

func TestToTeamResponseList_PreservesIDs(t *testing.T) {
	// Given: A slice with teams with specific IDs
	teams := []models.Team{
		{Model: gorm.Model{ID: 100}, Name: "Team 100"},
		{Model: gorm.Model{ID: 200}, Name: "Team 200"},
	}

	// When: Converting to response list
	responses := ToTeamResponseList(teams)

	// Then: IDs should be preserved
	assert.Equal(t, uint(100), responses[0].ID)
	assert.Equal(t, uint(200), responses[1].ID)
}
