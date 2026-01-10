package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/mocks"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestTeamHandler_GetAllTeams_Success_Unit(t *testing.T) {
	// Given: Teams exist in the repository
	mockTeamRepo := new(mocks.MockTeamRepository)
	handler := NewTeamHandlerWithRepo(mockTeamRepo)

	app := fiber.New()
	app.Get("/teams", handler.GetAllTeams)

	teams := []models.Team{
		{Model: gorm.Model{ID: 1}, Name: "Team Alpha"},
		{Model: gorm.Model{ID: 2}, Name: "Team Beta"},
	}
	mockTeamRepo.On("FindAll", mock.Anything).Return(teams, nil)

	req := httptest.NewRequest("GET", "/teams", nil)

	// When: Making the get all teams request
	resp, err := app.Test(req)

	// Then: The request should succeed with teams list
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockTeamRepo.AssertExpectations(t)
}

func TestTeamHandler_GetAllTeams_Empty_Unit(t *testing.T) {
	// Given: No teams exist in the repository
	mockTeamRepo := new(mocks.MockTeamRepository)
	handler := NewTeamHandlerWithRepo(mockTeamRepo)

	app := fiber.New()
	app.Get("/teams", handler.GetAllTeams)

	mockTeamRepo.On("FindAll", mock.Anything).Return([]models.Team{}, nil)

	req := httptest.NewRequest("GET", "/teams", nil)

	// When: Making the get all teams request
	resp, err := app.Test(req)

	// Then: The request should succeed with empty list
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockTeamRepo.AssertExpectations(t)
}

func TestTeamHandler_GetAllTeams_DatabaseError_Unit(t *testing.T) {
	// Given: A database error occurs
	mockTeamRepo := new(mocks.MockTeamRepository)
	handler := NewTeamHandlerWithRepo(mockTeamRepo)

	app := fiber.New()
	app.Get("/teams", handler.GetAllTeams)

	mockTeamRepo.On("FindAll", mock.Anything).Return(nil, errors.New("database error"))

	req := httptest.NewRequest("GET", "/teams", nil)

	// When: Making the get all teams request
	resp, err := app.Test(req)

	// Then: The request should fail with internal server error
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockTeamRepo.AssertExpectations(t)
}

func TestTeamHandler_GetTeamByID_Success_Unit(t *testing.T) {
	// Given: A team exists with the specified ID
	mockTeamRepo := new(mocks.MockTeamRepository)
	handler := NewTeamHandlerWithRepo(mockTeamRepo)

	app := fiber.New()
	app.Get("/teams/:id", handler.GetTeamByID)

	team := &models.Team{Model: gorm.Model{ID: 1}, Name: "Test Team"}
	mockTeamRepo.On("FindByID", mock.Anything, "1").Return(team, nil)

	req := httptest.NewRequest("GET", "/teams/1", nil)

	// When: Making the get team by ID request
	resp, err := app.Test(req)

	// Then: The request should succeed with team data
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockTeamRepo.AssertExpectations(t)
}

func TestTeamHandler_GetTeamByID_NotFound_Unit(t *testing.T) {
	// Given: No team exists with the specified ID
	mockTeamRepo := new(mocks.MockTeamRepository)
	handler := NewTeamHandlerWithRepo(mockTeamRepo)

	app := fiber.New()
	app.Get("/teams/:id", handler.GetTeamByID)

	mockTeamRepo.On("FindByID", mock.Anything, "999").Return(nil, gorm.ErrRecordNotFound)

	req := httptest.NewRequest("GET", "/teams/999", nil)

	// When: Making the get team by ID request
	resp, err := app.Test(req)

	// Then: The request should fail with not found status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockTeamRepo.AssertExpectations(t)
}

func TestTeamHandler_CreateTeam_Success_Unit(t *testing.T) {
	// Given: A valid create team request
	mockTeamRepo := new(mocks.MockTeamRepository)
	handler := NewTeamHandlerWithRepo(mockTeamRepo)

	app := fiber.New()
	app.Post("/teams", handler.CreateTeam)

	mockTeamRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Team")).Return(nil)

	reqBody := dtos.CreateTeamRequest{
		Name: "New Team",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/teams", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create team request
	resp, err := app.Test(req)

	// Then: The request should succeed with created status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	mockTeamRepo.AssertExpectations(t)
}

func TestTeamHandler_CreateTeam_InvalidJSON_Unit(t *testing.T) {
	// Given: An invalid JSON request body
	mockTeamRepo := new(mocks.MockTeamRepository)
	handler := NewTeamHandlerWithRepo(mockTeamRepo)

	app := fiber.New()
	app.Post("/teams", handler.CreateTeam)

	req := httptest.NewRequest("POST", "/teams", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create team request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestTeamHandler_CreateTeam_DatabaseError_Unit(t *testing.T) {
	// Given: A database error occurs when creating team
	mockTeamRepo := new(mocks.MockTeamRepository)
	handler := NewTeamHandlerWithRepo(mockTeamRepo)

	app := fiber.New()
	app.Post("/teams", handler.CreateTeam)

	mockTeamRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Team")).Return(errors.New("database error"))

	reqBody := dtos.CreateTeamRequest{
		Name: "New Team",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/teams", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create team request
	resp, err := app.Test(req)

	// Then: The request should fail with internal server error
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockTeamRepo.AssertExpectations(t)
}

func TestTeamHandler_UpdateTeam_Success_Unit(t *testing.T) {
	// Given: A team exists and valid update request
	mockTeamRepo := new(mocks.MockTeamRepository)
	handler := NewTeamHandlerWithRepo(mockTeamRepo)

	app := fiber.New()
	app.Put("/teams/:id", handler.UpdateTeam)

	existingTeam := &models.Team{Model: gorm.Model{ID: 1}, Name: "Old Name"}
	mockTeamRepo.On("FindByID", mock.Anything, "1").Return(existingTeam, nil)
	mockTeamRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Team")).Return(nil)

	reqBody := dtos.CreateTeamRequest{
		Name: "New Name",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/teams/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update team request
	resp, err := app.Test(req)

	// Then: The request should succeed
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockTeamRepo.AssertExpectations(t)
}

func TestTeamHandler_UpdateTeam_NotFound_Unit(t *testing.T) {
	// Given: No team exists with the specified ID
	mockTeamRepo := new(mocks.MockTeamRepository)
	handler := NewTeamHandlerWithRepo(mockTeamRepo)

	app := fiber.New()
	app.Put("/teams/:id", handler.UpdateTeam)

	mockTeamRepo.On("FindByID", mock.Anything, "999").Return(nil, gorm.ErrRecordNotFound)

	reqBody := dtos.CreateTeamRequest{
		Name: "New Name",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/teams/999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update team request
	resp, err := app.Test(req)

	// Then: The request should fail with not found status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockTeamRepo.AssertExpectations(t)
}

func TestTeamHandler_UpdateTeam_InvalidJSON_Unit(t *testing.T) {
	// Given: A team exists but request body is invalid JSON
	mockTeamRepo := new(mocks.MockTeamRepository)
	handler := NewTeamHandlerWithRepo(mockTeamRepo)

	app := fiber.New()
	app.Put("/teams/:id", handler.UpdateTeam)

	existingTeam := &models.Team{Model: gorm.Model{ID: 1}, Name: "Old Name"}
	mockTeamRepo.On("FindByID", mock.Anything, "1").Return(existingTeam, nil)

	req := httptest.NewRequest("PUT", "/teams/1", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update team request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	mockTeamRepo.AssertExpectations(t)
}

func TestTeamHandler_UpdateTeam_DatabaseError_Unit(t *testing.T) {
	// Given: A team exists but database error occurs during update
	mockTeamRepo := new(mocks.MockTeamRepository)
	handler := NewTeamHandlerWithRepo(mockTeamRepo)

	app := fiber.New()
	app.Put("/teams/:id", handler.UpdateTeam)

	existingTeam := &models.Team{Model: gorm.Model{ID: 1}, Name: "Old Name"}
	mockTeamRepo.On("FindByID", mock.Anything, "1").Return(existingTeam, nil)
	mockTeamRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Team")).Return(errors.New("database error"))

	reqBody := dtos.CreateTeamRequest{
		Name: "New Name",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/teams/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update team request
	resp, err := app.Test(req)

	// Then: The request should fail with internal server error
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockTeamRepo.AssertExpectations(t)
}

func TestTeamHandler_DeleteTeam_Success_Unit(t *testing.T) {
	// Given: A team exists with the specified ID
	mockTeamRepo := new(mocks.MockTeamRepository)
	handler := NewTeamHandlerWithRepo(mockTeamRepo)

	app := fiber.New()
	app.Delete("/teams/:id", handler.DeleteTeam)

	existingTeam := &models.Team{Model: gorm.Model{ID: 1}, Name: "Team to Delete"}
	mockTeamRepo.On("FindByID", mock.Anything, "1").Return(existingTeam, nil)
	mockTeamRepo.On("Delete", mock.Anything, uint(1)).Return(nil)

	req := httptest.NewRequest("DELETE", "/teams/1", nil)

	// When: Making the delete team request
	resp, err := app.Test(req)

	// Then: The request should succeed with no content status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	mockTeamRepo.AssertExpectations(t)
}

func TestTeamHandler_DeleteTeam_NotFound_Unit(t *testing.T) {
	// Given: No team exists with the specified ID
	mockTeamRepo := new(mocks.MockTeamRepository)
	handler := NewTeamHandlerWithRepo(mockTeamRepo)

	app := fiber.New()
	app.Delete("/teams/:id", handler.DeleteTeam)

	mockTeamRepo.On("FindByID", mock.Anything, "999").Return(nil, gorm.ErrRecordNotFound)

	req := httptest.NewRequest("DELETE", "/teams/999", nil)

	// When: Making the delete team request
	resp, err := app.Test(req)

	// Then: The request should fail with not found status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockTeamRepo.AssertExpectations(t)
}

func TestTeamHandler_DeleteTeam_DatabaseError_Unit(t *testing.T) {
	// Given: A team exists but database error occurs during delete
	mockTeamRepo := new(mocks.MockTeamRepository)
	handler := NewTeamHandlerWithRepo(mockTeamRepo)

	app := fiber.New()
	app.Delete("/teams/:id", handler.DeleteTeam)

	existingTeam := &models.Team{Model: gorm.Model{ID: 1}, Name: "Team to Delete"}
	mockTeamRepo.On("FindByID", mock.Anything, "1").Return(existingTeam, nil)
	mockTeamRepo.On("Delete", mock.Anything, uint(1)).Return(errors.New("database error"))

	req := httptest.NewRequest("DELETE", "/teams/1", nil)

	// When: Making the delete team request
	resp, err := app.Test(req)

	// Then: The request should fail with internal server error
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockTeamRepo.AssertExpectations(t)
}
