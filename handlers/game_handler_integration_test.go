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

func setupGameTestApp(db *gorm.DB) *fiber.App {
	app := fiber.New()
	gameHandler := NewGameHandler(db)

	app.Get("/games", gameHandler.GetAllGames)
	app.Get("/games/:id", gameHandler.GetGameByID)
	app.Post("/games", gameHandler.CreateGame)
	app.Put("/games/:id", gameHandler.UpdateGame)
	app.Delete("/games/:id", gameHandler.DeleteGame)

	return app
}

func TestGameHandler_GetAllGames_Empty(t *testing.T) {
	// Given: An empty database
	db := setupTestDB(t)
	app := setupGameTestApp(db)

	req := httptest.NewRequest("GET", "/games", nil)

	// When: Making the get all games request
	resp, err := app.Test(req)

	// Then: The request should succeed with empty array
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var games []dtos.GameResponse
	json.NewDecoder(resp.Body).Decode(&games)
	assert.Empty(t, games)
}

func TestGameHandler_GetAllGames_WithGames(t *testing.T) {
	// Given: A database with games
	db := setupTestDB(t)
	app := setupGameTestApp(db)

	db.Create(&models.Game{Name: "Game 1", Description: "First game"})
	db.Create(&models.Game{Name: "Game 2", Description: "Second game"})

	req := httptest.NewRequest("GET", "/games", nil)

	// When: Making the get all games request
	resp, err := app.Test(req)

	// Then: The request should return all games
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var games []dtos.GameResponse
	json.NewDecoder(resp.Body).Decode(&games)
	assert.Len(t, games, 2)
}

func TestGameHandler_GetGameByID_Found(t *testing.T) {
	// Given: A game exists in the database
	db := setupTestDB(t)
	app := setupGameTestApp(db)

	game := models.Game{Name: "Chess", Description: "Strategy game"}
	db.Create(&game)

	req := httptest.NewRequest("GET", fmt.Sprintf("/games/%d", game.ID), nil)

	// When: Making the get game by ID request
	resp, err := app.Test(req)

	// Then: The request should return the game
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response dtos.GameResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Chess", response.Name)
}

func TestGameHandler_GetGameByID_NotFound(t *testing.T) {
	// Given: An empty database
	db := setupTestDB(t)
	app := setupGameTestApp(db)

	req := httptest.NewRequest("GET", "/games/999", nil)

	// When: Making the get game by ID request for non-existent game
	resp, err := app.Test(req)

	// Then: The request should return not found
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestGameHandler_CreateGame_Success(t *testing.T) {
	// Given: A valid create game request
	db := setupTestDB(t)
	app := setupGameTestApp(db)

	reqBody := dtos.CreateGameRequest{
		Name:            "New Game",
		Description:     "A new board game",
		NumberOfPlayers: 4,
		MinPlayers:      2,
		MaxPlayers:      4,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/games", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create game request
	resp, err := app.Test(req)

	// Then: The game should be created
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestGameHandler_CreateGame_InvalidJSON(t *testing.T) {
	// Given: An invalid JSON request body
	db := setupTestDB(t)
	app := setupGameTestApp(db)

	req := httptest.NewRequest("POST", "/games", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create game request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestGameHandler_CreateGame_ReturnsCreatedGame(t *testing.T) {
	// Given: A valid create game request
	db := setupTestDB(t)
	app := setupGameTestApp(db)

	reqBody := dtos.CreateGameRequest{
		Name:            "Created Game",
		Description:     "Description",
		NumberOfPlayers: 2,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/games", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create game request
	resp, err := app.Test(req)

	// Then: The response should contain the created game
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var response dtos.GameResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Created Game", response.Name)
	assert.NotEqual(t, uint(0), response.ID)
}

func TestGameHandler_UpdateGame_Success(t *testing.T) {
	// Given: An existing game
	db := setupTestDB(t)
	app := setupGameTestApp(db)

	game := models.Game{Name: "Original", Description: "Original description"}
	db.Create(&game)

	reqBody := dtos.CreateGameRequest{
		Name:            "Updated",
		Description:     "Updated description",
		NumberOfPlayers: 4,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", fmt.Sprintf("/games/%d", game.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update game request
	resp, err := app.Test(req)

	// Then: The game should be updated
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response dtos.GameResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Updated", response.Name)
}

func TestGameHandler_UpdateGame_NotFound(t *testing.T) {
	// Given: An empty database
	db := setupTestDB(t)
	app := setupGameTestApp(db)

	reqBody := dtos.CreateGameRequest{Name: "Updated"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/games/999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update request for non-existent game
	resp, err := app.Test(req)

	// Then: The request should return not found
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestGameHandler_UpdateGame_InvalidJSON(t *testing.T) {
	// Given: An existing game and invalid JSON
	db := setupTestDB(t)
	app := setupGameTestApp(db)

	game := models.Game{Name: "Original"}
	db.Create(&game)

	req := httptest.NewRequest("PUT", fmt.Sprintf("/games/%d", game.ID), bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update request with invalid JSON
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestGameHandler_DeleteGame_Success(t *testing.T) {
	// Given: An existing game
	db := setupTestDB(t)
	app := setupGameTestApp(db)

	game := models.Game{Name: "To Delete"}
	db.Create(&game)

	req := httptest.NewRequest("DELETE", fmt.Sprintf("/games/%d", game.ID), nil)

	// When: Making the delete game request
	resp, err := app.Test(req)

	// Then: The game should be deleted
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
}

func TestGameHandler_DeleteGame_NotFound(t *testing.T) {
	// Given: An empty database
	db := setupTestDB(t)
	app := setupGameTestApp(db)

	req := httptest.NewRequest("DELETE", "/games/999", nil)

	// When: Making the delete request for non-existent game
	resp, err := app.Test(req)

	// Then: The request should return not found
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestGameHandler_DeleteGame_ActuallyDeletes(t *testing.T) {
	// Given: An existing game
	db := setupTestDB(t)
	app := setupGameTestApp(db)

	game := models.Game{Name: "To Delete"}
	db.Create(&game)
	gameID := game.ID

	req := httptest.NewRequest("DELETE", fmt.Sprintf("/games/%d", gameID), nil)

	// When: Making the delete game request
	app.Test(req)

	// Then: The game should no longer exist in the database
	var count int64
	db.Model(&models.Game{}).Where("id = ?", gameID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestNewGameHandler(t *testing.T) {
	// Given: A database connection
	db := setupTestDB(t)

	// When: Creating a new game handler
	handler := NewGameHandler(db)

	// Then: The handler should be created
	assert.NotNil(t, handler)
}
