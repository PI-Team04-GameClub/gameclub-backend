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

func TestGameHandler_GetAllGames_Success_Unit(t *testing.T) {
	// Given: Games exist in the repository
	mockGameRepo := new(mocks.MockGameRepository)
	handler := NewGameHandlerWithRepo(mockGameRepo)

	app := fiber.New()
	app.Get("/games", handler.GetAllGames)

	games := []models.Game{
		{Model: gorm.Model{ID: 1}, Name: "Game 1", Category: models.CategoryStrategy},
		{Model: gorm.Model{ID: 2}, Name: "Game 2", Category: models.CategoryParty},
	}
	mockGameRepo.On("FindAll", mock.Anything).Return(games, nil)

	req := httptest.NewRequest("GET", "/games", nil)

	// When: Making the get all games request
	resp, err := app.Test(req)

	// Then: The request should succeed with games list
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockGameRepo.AssertExpectations(t)
}

func TestGameHandler_GetAllGames_Empty_Unit(t *testing.T) {
	// Given: No games exist in the repository
	mockGameRepo := new(mocks.MockGameRepository)
	handler := NewGameHandlerWithRepo(mockGameRepo)

	app := fiber.New()
	app.Get("/games", handler.GetAllGames)

	mockGameRepo.On("FindAll", mock.Anything).Return([]models.Game{}, nil)

	req := httptest.NewRequest("GET", "/games", nil)

	// When: Making the get all games request
	resp, err := app.Test(req)

	// Then: The request should succeed with empty list
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockGameRepo.AssertExpectations(t)
}

func TestGameHandler_GetAllGames_DatabaseError_Unit(t *testing.T) {
	// Given: A database error occurs
	mockGameRepo := new(mocks.MockGameRepository)
	handler := NewGameHandlerWithRepo(mockGameRepo)

	app := fiber.New()
	app.Get("/games", handler.GetAllGames)

	mockGameRepo.On("FindAll", mock.Anything).Return(nil, errors.New("database error"))

	req := httptest.NewRequest("GET", "/games", nil)

	// When: Making the get all games request
	resp, err := app.Test(req)

	// Then: The request should fail with internal server error
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockGameRepo.AssertExpectations(t)
}

func TestGameHandler_GetGameByID_Success_Unit(t *testing.T) {
	// Given: A game exists with the specified ID
	mockGameRepo := new(mocks.MockGameRepository)
	handler := NewGameHandlerWithRepo(mockGameRepo)

	app := fiber.New()
	app.Get("/games/:id", handler.GetGameByID)

	game := &models.Game{Model: gorm.Model{ID: 1}, Name: "Test Game", Category: models.CategoryStrategy}
	mockGameRepo.On("FindByID", mock.Anything, "1").Return(game, nil)

	req := httptest.NewRequest("GET", "/games/1", nil)

	// When: Making the get game by ID request
	resp, err := app.Test(req)

	// Then: The request should succeed with game data
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockGameRepo.AssertExpectations(t)
}

func TestGameHandler_GetGameByID_NotFound_Unit(t *testing.T) {
	// Given: No game exists with the specified ID
	mockGameRepo := new(mocks.MockGameRepository)
	handler := NewGameHandlerWithRepo(mockGameRepo)

	app := fiber.New()
	app.Get("/games/:id", handler.GetGameByID)

	mockGameRepo.On("FindByID", mock.Anything, "999").Return(nil, gorm.ErrRecordNotFound)

	req := httptest.NewRequest("GET", "/games/999", nil)

	// When: Making the get game by ID request
	resp, err := app.Test(req)

	// Then: The request should fail with not found status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockGameRepo.AssertExpectations(t)
}

func TestGameHandler_CreateGame_Success_Unit(t *testing.T) {
	// Given: A valid create game request
	mockGameRepo := new(mocks.MockGameRepository)
	handler := NewGameHandlerWithRepo(mockGameRepo)

	app := fiber.New()
	app.Post("/games", handler.CreateGame)

	mockGameRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Game")).Return(nil)

	reqBody := dtos.CreateGameRequest{
		Name:            "New Game",
		NumberOfPlayers: 4,
		Category:        "Strategy",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/games", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create game request
	resp, err := app.Test(req)

	// Then: The request should succeed with created status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	mockGameRepo.AssertExpectations(t)
}

func TestGameHandler_CreateGame_InvalidJSON_Unit(t *testing.T) {
	// Given: An invalid JSON request body
	mockGameRepo := new(mocks.MockGameRepository)
	handler := NewGameHandlerWithRepo(mockGameRepo)

	app := fiber.New()
	app.Post("/games", handler.CreateGame)

	req := httptest.NewRequest("POST", "/games", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create game request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestGameHandler_CreateGame_DatabaseError_Unit(t *testing.T) {
	// Given: A database error occurs when creating game
	mockGameRepo := new(mocks.MockGameRepository)
	handler := NewGameHandlerWithRepo(mockGameRepo)

	app := fiber.New()
	app.Post("/games", handler.CreateGame)

	mockGameRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Game")).Return(errors.New("database error"))

	reqBody := dtos.CreateGameRequest{
		Name:            "New Game",
		NumberOfPlayers: 4,
		Category:        "Strategy",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/games", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create game request
	resp, err := app.Test(req)

	// Then: The request should fail with internal server error
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockGameRepo.AssertExpectations(t)
}

func TestGameHandler_UpdateGame_Success_Unit(t *testing.T) {
	// Given: A game exists and valid update request
	mockGameRepo := new(mocks.MockGameRepository)
	handler := NewGameHandlerWithRepo(mockGameRepo)

	app := fiber.New()
	app.Put("/games/:id", handler.UpdateGame)

	existingGame := &models.Game{Model: gorm.Model{ID: 1}, Name: "Old Name", Category: models.CategoryStrategy}
	mockGameRepo.On("FindByID", mock.Anything, "1").Return(existingGame, nil)
	mockGameRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Game")).Return(nil)

	reqBody := dtos.CreateGameRequest{
		Name:            "New Name",
		NumberOfPlayers: 4,
		Category:        "Party",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/games/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update game request
	resp, err := app.Test(req)

	// Then: The request should succeed
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockGameRepo.AssertExpectations(t)
}

func TestGameHandler_UpdateGame_NotFound_Unit(t *testing.T) {
	// Given: No game exists with the specified ID
	mockGameRepo := new(mocks.MockGameRepository)
	handler := NewGameHandlerWithRepo(mockGameRepo)

	app := fiber.New()
	app.Put("/games/:id", handler.UpdateGame)

	mockGameRepo.On("FindByID", mock.Anything, "999").Return(nil, gorm.ErrRecordNotFound)

	reqBody := dtos.CreateGameRequest{
		Name:            "New Name",
		NumberOfPlayers: 4,
		Category:        "Party",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/games/999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update game request
	resp, err := app.Test(req)

	// Then: The request should fail with not found status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockGameRepo.AssertExpectations(t)
}

func TestGameHandler_UpdateGame_InvalidJSON_Unit(t *testing.T) {
	// Given: A game exists but request body is invalid JSON
	mockGameRepo := new(mocks.MockGameRepository)
	handler := NewGameHandlerWithRepo(mockGameRepo)

	app := fiber.New()
	app.Put("/games/:id", handler.UpdateGame)

	existingGame := &models.Game{Model: gorm.Model{ID: 1}, Name: "Old Name", Category: models.CategoryStrategy}
	mockGameRepo.On("FindByID", mock.Anything, "1").Return(existingGame, nil)

	req := httptest.NewRequest("PUT", "/games/1", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update game request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	mockGameRepo.AssertExpectations(t)
}

func TestGameHandler_UpdateGame_DatabaseError_Unit(t *testing.T) {
	// Given: A game exists but database error occurs during update
	mockGameRepo := new(mocks.MockGameRepository)
	handler := NewGameHandlerWithRepo(mockGameRepo)

	app := fiber.New()
	app.Put("/games/:id", handler.UpdateGame)

	existingGame := &models.Game{Model: gorm.Model{ID: 1}, Name: "Old Name", Category: models.CategoryStrategy}
	mockGameRepo.On("FindByID", mock.Anything, "1").Return(existingGame, nil)
	mockGameRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Game")).Return(errors.New("database error"))

	reqBody := dtos.CreateGameRequest{
		Name:            "New Name",
		NumberOfPlayers: 4,
		Category:        "Party",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/games/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update game request
	resp, err := app.Test(req)

	// Then: The request should fail with internal server error
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockGameRepo.AssertExpectations(t)
}

func TestGameHandler_DeleteGame_Success_Unit(t *testing.T) {
	// Given: A game exists with the specified ID
	mockGameRepo := new(mocks.MockGameRepository)
	handler := NewGameHandlerWithRepo(mockGameRepo)

	app := fiber.New()
	app.Delete("/games/:id", handler.DeleteGame)

	existingGame := &models.Game{Model: gorm.Model{ID: 1}, Name: "Game to Delete", Category: models.CategoryStrategy}
	mockGameRepo.On("FindByID", mock.Anything, "1").Return(existingGame, nil)
	mockGameRepo.On("Delete", mock.Anything, uint(1)).Return(nil)

	req := httptest.NewRequest("DELETE", "/games/1", nil)

	// When: Making the delete game request
	resp, err := app.Test(req)

	// Then: The request should succeed with no content status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	mockGameRepo.AssertExpectations(t)
}

func TestGameHandler_DeleteGame_NotFound_Unit(t *testing.T) {
	// Given: No game exists with the specified ID
	mockGameRepo := new(mocks.MockGameRepository)
	handler := NewGameHandlerWithRepo(mockGameRepo)

	app := fiber.New()
	app.Delete("/games/:id", handler.DeleteGame)

	mockGameRepo.On("FindByID", mock.Anything, "999").Return(nil, gorm.ErrRecordNotFound)

	req := httptest.NewRequest("DELETE", "/games/999", nil)

	// When: Making the delete game request
	resp, err := app.Test(req)

	// Then: The request should fail with not found status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockGameRepo.AssertExpectations(t)
}

func TestGameHandler_DeleteGame_DatabaseError_Unit(t *testing.T) {
	// Given: A game exists but database error occurs during delete
	mockGameRepo := new(mocks.MockGameRepository)
	handler := NewGameHandlerWithRepo(mockGameRepo)

	app := fiber.New()
	app.Delete("/games/:id", handler.DeleteGame)

	existingGame := &models.Game{Model: gorm.Model{ID: 1}, Name: "Game to Delete", Category: models.CategoryStrategy}
	mockGameRepo.On("FindByID", mock.Anything, "1").Return(existingGame, nil)
	mockGameRepo.On("Delete", mock.Anything, uint(1)).Return(errors.New("database error"))

	req := httptest.NewRequest("DELETE", "/games/1", nil)

	// When: Making the delete game request
	resp, err := app.Test(req)

	// Then: The request should fail with internal server error
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockGameRepo.AssertExpectations(t)
}
