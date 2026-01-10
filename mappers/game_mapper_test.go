package mappers

import (
	"testing"

	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestToGameResponse_BasicFields(t *testing.T) {
	// Given: A game model with all fields populated
	game := &models.Game{
		Model:           gorm.Model{ID: 1},
		Name:            "Chess",
		Description:     "A classic strategy game",
		NumberOfPlayers: 2,
		MinPlayers:      2,
		MaxPlayers:      2,
		PlaytimeMinutes: 60,
		MinAge:          6,
		Complexity:      models.ComplexityMedium,
		Category:        models.CategoryStrategy,
		Publisher:       "Classic Games",
		YearPublished:   1850,
		Rating:          9.5,
	}

	// When: Converting to response
	response := ToGameResponse(game)

	// Then: All fields should be correctly mapped
	assert.Equal(t, uint(1), response.ID)
	assert.Equal(t, "Chess", response.Name)
	assert.Equal(t, "A classic strategy game", response.Description)
	assert.Equal(t, 2, response.NumberOfPlayers)
	assert.Equal(t, 2, response.MinPlayers)
	assert.Equal(t, 2, response.MaxPlayers)
	assert.Equal(t, 60, response.PlaytimeMinutes)
	assert.Equal(t, 6, response.MinAge)
	assert.Equal(t, "Medium", response.Complexity)
	assert.Equal(t, "Strategy", response.Category)
	assert.Equal(t, "Classic Games", response.Publisher)
	assert.Equal(t, 1850, response.YearPublished)
	assert.Equal(t, 9.5, response.Rating)
}

func TestToGameResponse_EmptyFields(t *testing.T) {
	// Given: A game model with empty/zero fields
	game := &models.Game{
		Model: gorm.Model{ID: 2},
		Name:  "Minimal Game",
	}

	// When: Converting to response
	response := ToGameResponse(game)

	// Then: Empty fields should be preserved
	assert.Equal(t, uint(2), response.ID)
	assert.Equal(t, "Minimal Game", response.Name)
	assert.Equal(t, "", response.Description)
	assert.Equal(t, 0, response.NumberOfPlayers)
}

func TestToGameModel_WithPlayerRange(t *testing.T) {
	// Given: A create request with min and max players
	req := dtos.CreateGameRequest{
		Name:        "Party Game",
		Description: "Fun for everyone",
		MinPlayers:  3,
		MaxPlayers:  8,
	}

	// When: Converting to model
	game := ToGameModel(req)

	// Then: The player range should be set correctly
	assert.Equal(t, "Party Game", game.Name)
	assert.Equal(t, 3, game.MinPlayers)
	assert.Equal(t, 8, game.MaxPlayers)
	assert.Equal(t, 8, game.NumberOfPlayers)
}

func TestToGameModel_WithNumberOfPlayers(t *testing.T) {
	// Given: A create request with only numberOfPlayers
	req := dtos.CreateGameRequest{
		Name:            "Card Game",
		NumberOfPlayers: 4,
	}

	// When: Converting to model
	game := ToGameModel(req)

	// Then: MinPlayers should default to 2 and MaxPlayers should match numberOfPlayers
	assert.Equal(t, 4, game.NumberOfPlayers)
	assert.Equal(t, 2, game.MinPlayers)
	assert.Equal(t, 4, game.MaxPlayers)
}

func TestToGameModel_WithOptionalFields(t *testing.T) {
	// Given: A create request with optional fields
	req := dtos.CreateGameRequest{
		Name:            "Complex Game",
		PlaytimeMinutes: 120,
		MinAge:          14,
		Complexity:      "Hard",
		Category:        "Strategy",
		Publisher:       "Board Game Publisher",
		YearPublished:   2023,
		Rating:          8.5,
	}

	// When: Converting to model
	game := ToGameModel(req)

	// Then: Optional fields should be set correctly
	assert.Equal(t, 120, game.PlaytimeMinutes)
	assert.Equal(t, 14, game.MinAge)
	assert.Equal(t, models.GameComplexity("Hard"), game.Complexity)
	assert.Equal(t, models.GameCategory("Strategy"), game.Category)
	assert.Equal(t, "Board Game Publisher", game.Publisher)
	assert.Equal(t, 2023, game.YearPublished)
	assert.Equal(t, 8.5, game.Rating)
}

func TestToGameModel_DefaultValues(t *testing.T) {
	// Given: A create request with minimal fields
	req := dtos.CreateGameRequest{
		Name: "Default Game",
	}

	// When: Converting to model
	game := ToGameModel(req)

	// Then: Default values from GameBuilder should be used
	assert.Equal(t, 2, game.MinPlayers)
	assert.Equal(t, 4, game.MaxPlayers)
	assert.Equal(t, 30, game.PlaytimeMinutes)
	assert.Equal(t, 8, game.MinAge)
	assert.Equal(t, models.ComplexityMedium, game.Complexity)
	assert.Equal(t, models.CategoryStrategy, game.Category)
}

func TestToGameResponseList_EmptySlice(t *testing.T) {
	// Given: An empty slice of games
	games := []models.Game{}

	// When: Converting to response list
	responses := ToGameResponseList(games)

	// Then: The result should be an empty slice
	assert.Empty(t, responses)
	assert.Len(t, responses, 0)
}

func TestToGameResponseList_MultipleGames(t *testing.T) {
	// Given: A slice with multiple games
	games := []models.Game{
		{Model: gorm.Model{ID: 1}, Name: "Game 1"},
		{Model: gorm.Model{ID: 2}, Name: "Game 2"},
		{Model: gorm.Model{ID: 3}, Name: "Game 3"},
	}

	// When: Converting to response list
	responses := ToGameResponseList(games)

	// Then: All games should be converted
	assert.Len(t, responses, 3)
	assert.Equal(t, "Game 1", responses[0].Name)
	assert.Equal(t, "Game 2", responses[1].Name)
	assert.Equal(t, "Game 3", responses[2].Name)
}

func TestUpdateGameFromRequest_UpdatesAllFields(t *testing.T) {
	// Given: An existing game and an update request
	existingGame := &models.Game{
		Model:           gorm.Model{ID: 10},
		Name:            "Old Name",
		Description:     "Old Description",
		NumberOfPlayers: 2,
	}
	req := dtos.CreateGameRequest{
		Name:            "New Name",
		Description:     "New Description",
		MinPlayers:      2,
		MaxPlayers:      6,
		PlaytimeMinutes: 90,
		MinAge:          12,
		Complexity:      "Expert",
		Category:        "Party",
		Publisher:       "New Publisher",
		YearPublished:   2024,
		Rating:          9.0,
	}

	// When: Updating the game
	updatedGame := UpdateGameFromRequest(existingGame, req)

	// Then: The game should have updated values and retain the original ID
	assert.Equal(t, uint(10), updatedGame.ID)
	assert.Equal(t, "New Name", updatedGame.Name)
	assert.Equal(t, "New Description", updatedGame.Description)
	assert.Equal(t, 2, updatedGame.MinPlayers)
	assert.Equal(t, 6, updatedGame.MaxPlayers)
	assert.Equal(t, 90, updatedGame.PlaytimeMinutes)
	assert.Equal(t, 12, updatedGame.MinAge)
}

func TestUpdateGameFromRequest_PreservesID(t *testing.T) {
	// Given: An existing game with a specific ID
	existingGame := &models.Game{
		Model: gorm.Model{ID: 999},
		Name:  "Original",
	}
	req := dtos.CreateGameRequest{
		Name: "Updated",
	}

	// When: Updating the game
	updatedGame := UpdateGameFromRequest(existingGame, req)

	// Then: The ID should be preserved
	assert.Equal(t, uint(999), updatedGame.ID)
}

func TestToGameResponse_ComplexityEnumConversion(t *testing.T) {
	// Given: A game with Expert complexity
	game := &models.Game{
		Model:      gorm.Model{ID: 1},
		Name:       "Expert Game",
		Complexity: models.ComplexityExpert,
	}

	// When: Converting to response
	response := ToGameResponse(game)

	// Then: Complexity should be converted to string "Expert"
	assert.Equal(t, "Expert", response.Complexity)
}

func TestToGameResponse_CategoryEnumConversion(t *testing.T) {
	// Given: A game with Cooperative category
	game := &models.Game{
		Model:    gorm.Model{ID: 1},
		Name:     "Coop Game",
		Category: models.CategoryCooperative,
	}

	// When: Converting to response
	response := ToGameResponse(game)

	// Then: Category should be converted to string "Cooperative"
	assert.Equal(t, "Cooperative", response.Category)
}
