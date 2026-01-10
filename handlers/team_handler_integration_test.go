package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTeamTestApp(db *gorm.DB) *fiber.App {
	app := fiber.New()
	teamHandler := NewTeamHandler(db)

	app.Get("/teams", teamHandler.GetAllTeams)
	app.Get("/teams/:id", teamHandler.GetTeamByID)
	app.Post("/teams", teamHandler.CreateTeam)
	app.Put("/teams/:id", teamHandler.UpdateTeam)
	app.Delete("/teams/:id", teamHandler.DeleteTeam)

	return app
}

func TestTeamHandler_GetAllTeams_Empty(t *testing.T) {
	// Given: An empty database
	db := setupTestDB(t)
	app := setupTeamTestApp(db)

	req := httptest.NewRequest("GET", "/teams", nil)

	// When: Making the get all teams request
	resp, err := app.Test(req)

	// Then: The request should succeed with empty array
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var teams []dtos.TeamResponse
	json.NewDecoder(resp.Body).Decode(&teams)
	assert.Empty(t, teams)
}

func TestTeamHandler_GetAllTeams_WithTeams(t *testing.T) {
	// Given: A database with teams
	db := setupTestDB(t)
	app := setupTeamTestApp(db)

	db.Create(&models.Team{Name: "Team Alpha"})
	db.Create(&models.Team{Name: "Team Beta"})
	db.Create(&models.Team{Name: "Team Gamma"})

	req := httptest.NewRequest("GET", "/teams", nil)

	// When: Making the get all teams request
	resp, err := app.Test(req)

	// Then: The request should return all teams
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var teams []dtos.TeamResponse
	json.NewDecoder(resp.Body).Decode(&teams)
	assert.Len(t, teams, 3)
}

func TestTeamHandler_GetTeamByID_Found(t *testing.T) {
	// Given: A team exists in the database
	db := setupTestDB(t)
	app := setupTeamTestApp(db)

	team := models.Team{Name: "Found Team"}
	db.Create(&team)

	req := httptest.NewRequest("GET", fmt.Sprintf("/teams/%d", team.ID), nil)

	// When: Making the get team by ID request
	resp, err := app.Test(req)

	// Then: The request should return the team
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response dtos.TeamResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Found Team", response.Name)
}

func TestTeamHandler_GetTeamByID_NotFound(t *testing.T) {
	// Given: An empty database
	db := setupTestDB(t)
	app := setupTeamTestApp(db)

	req := httptest.NewRequest("GET", "/teams/999", nil)

	// When: Making the get team by ID request for non-existent team
	resp, err := app.Test(req)

	// Then: The request should return not found
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestTeamHandler_CreateTeam_Success(t *testing.T) {
	// Given: A valid create team request
	db := setupTestDB(t)
	app := setupTeamTestApp(db)

	reqBody := dtos.CreateTeamRequest{
		Name: "New Team",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/teams", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create team request
	resp, err := app.Test(req)

	// Then: The team should be created
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestTeamHandler_CreateTeam_InvalidJSON(t *testing.T) {
	// Given: An invalid JSON request body
	db := setupTestDB(t)
	app := setupTeamTestApp(db)

	req := httptest.NewRequest("POST", "/teams", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create team request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestTeamHandler_CreateTeam_ReturnsCreatedTeam(t *testing.T) {
	// Given: A valid create team request
	db := setupTestDB(t)
	app := setupTeamTestApp(db)

	reqBody := dtos.CreateTeamRequest{
		Name: "Created Team",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/teams", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create team request
	resp, err := app.Test(req)

	// Then: The response should contain the created team
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var response dtos.TeamResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Created Team", response.Name)
	assert.NotEqual(t, uint(0), response.ID)
}

func TestTeamHandler_UpdateTeam_Success(t *testing.T) {
	// Given: An existing team
	db := setupTestDB(t)
	app := setupTeamTestApp(db)

	team := models.Team{Name: "Original Team"}
	db.Create(&team)

	reqBody := dtos.CreateTeamRequest{
		Name: "Updated Team",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", fmt.Sprintf("/teams/%d", team.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update team request
	resp, err := app.Test(req)

	// Then: The team should be updated
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response dtos.TeamResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Updated Team", response.Name)
}

func TestTeamHandler_UpdateTeam_NotFound(t *testing.T) {
	// Given: An empty database
	db := setupTestDB(t)
	app := setupTeamTestApp(db)

	reqBody := dtos.CreateTeamRequest{Name: "Updated"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/teams/999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update request for non-existent team
	resp, err := app.Test(req)

	// Then: The request should return not found
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestTeamHandler_UpdateTeam_InvalidJSON(t *testing.T) {
	// Given: An existing team and invalid JSON
	db := setupTestDB(t)
	app := setupTeamTestApp(db)

	team := models.Team{Name: "Original"}
	db.Create(&team)

	req := httptest.NewRequest("PUT", fmt.Sprintf("/teams/%d", team.ID), bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update request with invalid JSON
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestTeamHandler_DeleteTeam_Success(t *testing.T) {
	// Given: An existing team
	db := setupTestDB(t)
	app := setupTeamTestApp(db)

	team := models.Team{Name: "To Delete"}
	db.Create(&team)

	req := httptest.NewRequest("DELETE", fmt.Sprintf("/teams/%d", team.ID), nil)

	// When: Making the delete team request
	resp, err := app.Test(req)

	// Then: The team should be deleted
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
}

func TestTeamHandler_DeleteTeam_NotFound(t *testing.T) {
	// Given: An empty database
	db := setupTestDB(t)
	app := setupTeamTestApp(db)

	req := httptest.NewRequest("DELETE", "/teams/999", nil)

	// When: Making the delete request for non-existent team
	resp, err := app.Test(req)

	// Then: The request should return not found
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestTeamHandler_DeleteTeam_ActuallyDeletes(t *testing.T) {
	// Given: An existing team
	db := setupTestDB(t)
	app := setupTeamTestApp(db)

	team := models.Team{Name: "To Delete"}
	db.Create(&team)
	teamID := team.ID

	req := httptest.NewRequest("DELETE", fmt.Sprintf("/teams/%d", teamID), nil)

	// When: Making the delete team request
	app.Test(req)

	// Then: The team should no longer exist in the database
	var count int64
	db.Model(&models.Team{}).Where("id = ?", teamID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestNewTeamHandler(t *testing.T) {
	// Given: A database connection
	db := setupTestDB(t)

	// When: Creating a new team handler
	handler := NewTeamHandler(db)

	// Then: The handler should be created
	assert.NotNil(t, handler)
}
