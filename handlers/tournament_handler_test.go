package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/mocks"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestTournamentHandler_GetTournaments_Success_Unit(t *testing.T) {
	// Given: Tournaments exist in the repository
	mockTournamentRepo := new(mocks.MockTournamentRepository)
	mockGameRepo := new(mocks.MockGameRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewTournamentHandlerWithRepo(mockTournamentRepo, mockGameRepo, mockUserRepo)

	app := fiber.New()
	app.Get("/tournaments", handler.GetTournaments)

	tournaments := []models.Tournament{
		{Model: gorm.Model{ID: 1}, Name: "Tournament 1"},
		{Model: gorm.Model{ID: 2}, Name: "Tournament 2"},
	}
	mockTournamentRepo.On("FindAll", mock.Anything).Return(tournaments, nil)

	req := httptest.NewRequest("GET", "/tournaments", nil)

	// When: Making the get all tournaments request
	resp, err := app.Test(req)

	// Then: The request should succeed with tournaments list
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockTournamentRepo.AssertExpectations(t)
}

func TestTournamentHandler_GetTournaments_Empty_Unit(t *testing.T) {
	// Given: No tournaments exist in the repository
	mockTournamentRepo := new(mocks.MockTournamentRepository)
	mockGameRepo := new(mocks.MockGameRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewTournamentHandlerWithRepo(mockTournamentRepo, mockGameRepo, mockUserRepo)

	app := fiber.New()
	app.Get("/tournaments", handler.GetTournaments)

	mockTournamentRepo.On("FindAll", mock.Anything).Return([]models.Tournament{}, nil)

	req := httptest.NewRequest("GET", "/tournaments", nil)

	// When: Making the get all tournaments request
	resp, err := app.Test(req)

	// Then: The request should succeed with empty list
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockTournamentRepo.AssertExpectations(t)
}

func TestTournamentHandler_GetTournaments_DatabaseError_Unit(t *testing.T) {
	// Given: A database error occurs
	mockTournamentRepo := new(mocks.MockTournamentRepository)
	mockGameRepo := new(mocks.MockGameRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewTournamentHandlerWithRepo(mockTournamentRepo, mockGameRepo, mockUserRepo)

	app := fiber.New()
	app.Get("/tournaments", handler.GetTournaments)

	mockTournamentRepo.On("FindAll", mock.Anything).Return(nil, errors.New("database error"))

	req := httptest.NewRequest("GET", "/tournaments", nil)

	// When: Making the get all tournaments request
	resp, err := app.Test(req)

	// Then: The request should fail with internal server error
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockTournamentRepo.AssertExpectations(t)
}

func TestTournamentHandler_GetTournamentByID_Success_Unit(t *testing.T) {
	// Given: A tournament exists with the specified ID
	mockTournamentRepo := new(mocks.MockTournamentRepository)
	mockGameRepo := new(mocks.MockGameRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewTournamentHandlerWithRepo(mockTournamentRepo, mockGameRepo, mockUserRepo)

	app := fiber.New()
	app.Get("/tournaments/:id", handler.GetTournamentByID)

	tournament := &models.Tournament{Model: gorm.Model{ID: 1}, Name: "Test Tournament"}
	mockTournamentRepo.On("FindByID", mock.Anything, 1).Return(tournament, nil)

	req := httptest.NewRequest("GET", "/tournaments/1", nil)

	// When: Making the get tournament by ID request
	resp, err := app.Test(req)

	// Then: The request should succeed with tournament data
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockTournamentRepo.AssertExpectations(t)
}

func TestTournamentHandler_GetTournamentByID_NotFound_Unit(t *testing.T) {
	// Given: No tournament exists with the specified ID
	mockTournamentRepo := new(mocks.MockTournamentRepository)
	mockGameRepo := new(mocks.MockGameRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewTournamentHandlerWithRepo(mockTournamentRepo, mockGameRepo, mockUserRepo)

	app := fiber.New()
	app.Get("/tournaments/:id", handler.GetTournamentByID)

	mockTournamentRepo.On("FindByID", mock.Anything, 999).Return(nil, gorm.ErrRecordNotFound)

	req := httptest.NewRequest("GET", "/tournaments/999", nil)

	// When: Making the get tournament by ID request
	resp, err := app.Test(req)

	// Then: The request should fail with not found status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockTournamentRepo.AssertExpectations(t)
}

func TestTournamentHandler_GetTournamentByID_InvalidID_Unit(t *testing.T) {
	// Given: An invalid tournament ID
	mockTournamentRepo := new(mocks.MockTournamentRepository)
	mockGameRepo := new(mocks.MockGameRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewTournamentHandlerWithRepo(mockTournamentRepo, mockGameRepo, mockUserRepo)

	app := fiber.New()
	app.Get("/tournaments/:id", handler.GetTournamentByID)

	req := httptest.NewRequest("GET", "/tournaments/invalid", nil)

	// When: Making the get tournament by ID request with invalid ID
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestTournamentHandler_CreateTournament_Success_Unit(t *testing.T) {
	// Given: A valid create tournament request
	mockTournamentRepo := new(mocks.MockTournamentRepository)
	mockGameRepo := new(mocks.MockGameRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewTournamentHandlerWithRepo(mockTournamentRepo, mockGameRepo, mockUserRepo)

	app := fiber.New()
	app.Post("/tournaments", handler.CreateTournament)

	game := &models.Game{Model: gorm.Model{ID: 1}, Name: "Test Game"}
	mockGameRepo.On("FindByID", mock.Anything, "1").Return(game, nil)
	mockTournamentRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Tournament")).Return(nil)
	mockTournamentRepo.On("FindByID", mock.Anything, 0).Return(&models.Tournament{
		Model: gorm.Model{ID: 1},
		Name:  "New Tournament",
		Game:  *game,
	}, nil)
	mockUserRepo.On("FindAll", mock.Anything).Return([]models.User{}, nil)

	reqBody := dtos.CreateTournamentRequest{
		Name:      "New Tournament",
		GameId:    1,
		StartDate: time.Now(),
		PrizePool: 1000,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/tournaments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create tournament request
	resp, err := app.Test(req)

	// Then: The request should succeed with created status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	mockTournamentRepo.AssertExpectations(t)
	mockGameRepo.AssertExpectations(t)
}

func TestTournamentHandler_CreateTournament_InvalidJSON_Unit(t *testing.T) {
	// Given: An invalid JSON request body
	mockTournamentRepo := new(mocks.MockTournamentRepository)
	mockGameRepo := new(mocks.MockGameRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewTournamentHandlerWithRepo(mockTournamentRepo, mockGameRepo, mockUserRepo)

	app := fiber.New()
	app.Post("/tournaments", handler.CreateTournament)

	req := httptest.NewRequest("POST", "/tournaments", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create tournament request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestTournamentHandler_CreateTournament_GameNotFound_Unit(t *testing.T) {
	// Given: The referenced game does not exist
	mockTournamentRepo := new(mocks.MockTournamentRepository)
	mockGameRepo := new(mocks.MockGameRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewTournamentHandlerWithRepo(mockTournamentRepo, mockGameRepo, mockUserRepo)

	app := fiber.New()
	app.Post("/tournaments", handler.CreateTournament)

	mockGameRepo.On("FindByID", mock.Anything, "999").Return(nil, gorm.ErrRecordNotFound)

	reqBody := dtos.CreateTournamentRequest{
		Name:      "New Tournament",
		GameId:    999,
		StartDate: time.Now(),
		PrizePool: 1000,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/tournaments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create tournament request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	mockGameRepo.AssertExpectations(t)
}

func TestTournamentHandler_CreateTournament_DatabaseError_Unit(t *testing.T) {
	// Given: A database error occurs when creating tournament
	mockTournamentRepo := new(mocks.MockTournamentRepository)
	mockGameRepo := new(mocks.MockGameRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewTournamentHandlerWithRepo(mockTournamentRepo, mockGameRepo, mockUserRepo)

	app := fiber.New()
	app.Post("/tournaments", handler.CreateTournament)

	game := &models.Game{Model: gorm.Model{ID: 1}, Name: "Test Game"}
	mockGameRepo.On("FindByID", mock.Anything, "1").Return(game, nil)
	mockTournamentRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Tournament")).Return(errors.New("database error"))

	reqBody := dtos.CreateTournamentRequest{
		Name:      "New Tournament",
		GameId:    1,
		StartDate: time.Now(),
		PrizePool: 1000,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/tournaments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create tournament request
	resp, err := app.Test(req)

	// Then: The request should fail with internal server error
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockTournamentRepo.AssertExpectations(t)
	mockGameRepo.AssertExpectations(t)
}

func TestTournamentHandler_UpdateTournament_Success_Unit(t *testing.T) {
	// Given: A tournament exists and valid update request
	mockTournamentRepo := new(mocks.MockTournamentRepository)
	mockGameRepo := new(mocks.MockGameRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewTournamentHandlerWithRepo(mockTournamentRepo, mockGameRepo, mockUserRepo)

	app := fiber.New()
	app.Put("/tournaments/:id", handler.UpdateTournament)

	existingTournament := &models.Tournament{Model: gorm.Model{ID: 1}, Name: "Old Name"}
	game := &models.Game{Model: gorm.Model{ID: 1}, Name: "Test Game"}

	mockTournamentRepo.On("FindByID", mock.Anything, 1).Return(existingTournament, nil).Once()
	mockGameRepo.On("FindByID", mock.Anything, "1").Return(game, nil)
	mockTournamentRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Tournament")).Return(nil)
	mockTournamentRepo.On("FindByID", mock.Anything, 1).Return(&models.Tournament{
		Model: gorm.Model{ID: 1},
		Name:  "New Name",
		Game:  *game,
	}, nil).Once()

	reqBody := dtos.CreateTournamentRequest{
		Name:      "New Name",
		GameId:    1,
		StartDate: time.Now(),
		PrizePool: 2000,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/tournaments/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update tournament request
	resp, err := app.Test(req)

	// Then: The request should succeed
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockTournamentRepo.AssertExpectations(t)
	mockGameRepo.AssertExpectations(t)
}

func TestTournamentHandler_UpdateTournament_NotFound_Unit(t *testing.T) {
	// Given: No tournament exists with the specified ID
	mockTournamentRepo := new(mocks.MockTournamentRepository)
	mockGameRepo := new(mocks.MockGameRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewTournamentHandlerWithRepo(mockTournamentRepo, mockGameRepo, mockUserRepo)

	app := fiber.New()
	app.Put("/tournaments/:id", handler.UpdateTournament)

	mockTournamentRepo.On("FindByID", mock.Anything, 999).Return(nil, gorm.ErrRecordNotFound)

	reqBody := dtos.CreateTournamentRequest{
		Name:      "New Name",
		GameId:    1,
		StartDate: time.Now(),
		PrizePool: 2000,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/tournaments/999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update tournament request
	resp, err := app.Test(req)

	// Then: The request should fail with not found status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockTournamentRepo.AssertExpectations(t)
}

func TestTournamentHandler_UpdateTournament_InvalidID_Unit(t *testing.T) {
	// Given: An invalid tournament ID
	mockTournamentRepo := new(mocks.MockTournamentRepository)
	mockGameRepo := new(mocks.MockGameRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewTournamentHandlerWithRepo(mockTournamentRepo, mockGameRepo, mockUserRepo)

	app := fiber.New()
	app.Put("/tournaments/:id", handler.UpdateTournament)

	reqBody := dtos.CreateTournamentRequest{
		Name:   "New Name",
		GameId: 1,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/tournaments/invalid", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update tournament request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestTournamentHandler_UpdateTournament_InvalidJSON_Unit(t *testing.T) {
	// Given: A tournament exists but request body is invalid JSON
	mockTournamentRepo := new(mocks.MockTournamentRepository)
	mockGameRepo := new(mocks.MockGameRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewTournamentHandlerWithRepo(mockTournamentRepo, mockGameRepo, mockUserRepo)

	app := fiber.New()
	app.Put("/tournaments/:id", handler.UpdateTournament)

	existingTournament := &models.Tournament{Model: gorm.Model{ID: 1}, Name: "Old Name"}
	mockTournamentRepo.On("FindByID", mock.Anything, 1).Return(existingTournament, nil)

	req := httptest.NewRequest("PUT", "/tournaments/1", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update tournament request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	mockTournamentRepo.AssertExpectations(t)
}

func TestTournamentHandler_UpdateTournament_GameNotFound_Unit(t *testing.T) {
	// Given: The referenced game does not exist
	mockTournamentRepo := new(mocks.MockTournamentRepository)
	mockGameRepo := new(mocks.MockGameRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewTournamentHandlerWithRepo(mockTournamentRepo, mockGameRepo, mockUserRepo)

	app := fiber.New()
	app.Put("/tournaments/:id", handler.UpdateTournament)

	existingTournament := &models.Tournament{Model: gorm.Model{ID: 1}, Name: "Old Name"}
	mockTournamentRepo.On("FindByID", mock.Anything, 1).Return(existingTournament, nil)
	mockGameRepo.On("FindByID", mock.Anything, "999").Return(nil, gorm.ErrRecordNotFound)

	reqBody := dtos.CreateTournamentRequest{
		Name:   "New Name",
		GameId: 999,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/tournaments/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update tournament request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	mockTournamentRepo.AssertExpectations(t)
	mockGameRepo.AssertExpectations(t)
}

func TestTournamentHandler_UpdateTournament_DatabaseError_Unit(t *testing.T) {
	// Given: A tournament exists but database error occurs during update
	mockTournamentRepo := new(mocks.MockTournamentRepository)
	mockGameRepo := new(mocks.MockGameRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewTournamentHandlerWithRepo(mockTournamentRepo, mockGameRepo, mockUserRepo)

	app := fiber.New()
	app.Put("/tournaments/:id", handler.UpdateTournament)

	existingTournament := &models.Tournament{Model: gorm.Model{ID: 1}, Name: "Old Name"}
	game := &models.Game{Model: gorm.Model{ID: 1}, Name: "Test Game"}

	mockTournamentRepo.On("FindByID", mock.Anything, 1).Return(existingTournament, nil)
	mockGameRepo.On("FindByID", mock.Anything, "1").Return(game, nil)
	mockTournamentRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Tournament")).Return(errors.New("database error"))

	reqBody := dtos.CreateTournamentRequest{
		Name:   "New Name",
		GameId: 1,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/tournaments/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update tournament request
	resp, err := app.Test(req)

	// Then: The request should fail with internal server error
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockTournamentRepo.AssertExpectations(t)
	mockGameRepo.AssertExpectations(t)
}

func TestTournamentHandler_DeleteTournament_Success_Unit(t *testing.T) {
	// Given: A tournament exists with the specified ID
	mockTournamentRepo := new(mocks.MockTournamentRepository)
	mockGameRepo := new(mocks.MockGameRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewTournamentHandlerWithRepo(mockTournamentRepo, mockGameRepo, mockUserRepo)

	app := fiber.New()
	app.Delete("/tournaments/:id", handler.DeleteTournament)

	existingTournament := &models.Tournament{Model: gorm.Model{ID: 1}, Name: "Tournament to Delete"}
	mockTournamentRepo.On("FindByID", mock.Anything, 1).Return(existingTournament, nil)
	mockTournamentRepo.On("Delete", mock.Anything, 1).Return(nil)

	req := httptest.NewRequest("DELETE", "/tournaments/1", nil)

	// When: Making the delete tournament request
	resp, err := app.Test(req)

	// Then: The request should succeed with no content status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	mockTournamentRepo.AssertExpectations(t)
}

func TestTournamentHandler_DeleteTournament_NotFound_Unit(t *testing.T) {
	// Given: No tournament exists with the specified ID
	mockTournamentRepo := new(mocks.MockTournamentRepository)
	mockGameRepo := new(mocks.MockGameRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewTournamentHandlerWithRepo(mockTournamentRepo, mockGameRepo, mockUserRepo)

	app := fiber.New()
	app.Delete("/tournaments/:id", handler.DeleteTournament)

	mockTournamentRepo.On("FindByID", mock.Anything, 999).Return(nil, gorm.ErrRecordNotFound)

	req := httptest.NewRequest("DELETE", "/tournaments/999", nil)

	// When: Making the delete tournament request
	resp, err := app.Test(req)

	// Then: The request should fail with not found status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockTournamentRepo.AssertExpectations(t)
}

func TestTournamentHandler_DeleteTournament_InvalidID_Unit(t *testing.T) {
	// Given: An invalid tournament ID
	mockTournamentRepo := new(mocks.MockTournamentRepository)
	mockGameRepo := new(mocks.MockGameRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewTournamentHandlerWithRepo(mockTournamentRepo, mockGameRepo, mockUserRepo)

	app := fiber.New()
	app.Delete("/tournaments/:id", handler.DeleteTournament)

	req := httptest.NewRequest("DELETE", "/tournaments/invalid", nil)

	// When: Making the delete tournament request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestTournamentHandler_DeleteTournament_DatabaseError_Unit(t *testing.T) {
	// Given: A tournament exists but database error occurs during delete
	mockTournamentRepo := new(mocks.MockTournamentRepository)
	mockGameRepo := new(mocks.MockGameRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewTournamentHandlerWithRepo(mockTournamentRepo, mockGameRepo, mockUserRepo)

	app := fiber.New()
	app.Delete("/tournaments/:id", handler.DeleteTournament)

	existingTournament := &models.Tournament{Model: gorm.Model{ID: 1}, Name: "Tournament to Delete"}
	mockTournamentRepo.On("FindByID", mock.Anything, 1).Return(existingTournament, nil)
	mockTournamentRepo.On("Delete", mock.Anything, 1).Return(errors.New("database error"))

	req := httptest.NewRequest("DELETE", "/tournaments/1", nil)

	// When: Making the delete tournament request
	resp, err := app.Test(req)

	// Then: The request should fail with internal server error
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockTournamentRepo.AssertExpectations(t)
}
