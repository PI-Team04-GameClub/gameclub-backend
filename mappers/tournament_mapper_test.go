package mappers

import (
	"testing"
	"time"

	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestToTournamentResponse_BasicFields(t *testing.T) {
	// Given: A tournament model with all fields populated
	tournament := &models.Tournament{
		Model:               gorm.Model{ID: 1},
		Name:                "Summer Championship",
		GameID:              5,
		BasePrizePool:       1000.00,
		CalculatedPrizePool: 1200.00,
		BonusType:           "Summer Bonus (20%)",
		StartDate:           time.Date(2024, 7, 15, 10, 0, 0, 0, time.UTC),
		Status:              models.StatusUpcoming,
		Game:                models.Game{Name: "Chess"},
	}

	// When: Converting to response
	response := ToTournamentResponse(tournament)

	// Then: All fields should be correctly mapped
	assert.Equal(t, uint(1), response.ID)
	assert.Equal(t, "Summer Championship", response.Name)
	assert.Equal(t, "Chess", response.Game)
	assert.Equal(t, 1000.00, response.BasePrizePool)
	assert.Equal(t, 1200.00, response.CalculatedPrizePool)
	assert.Equal(t, 1200.00, response.PrizePool)
	assert.Equal(t, "Summer Bonus (20%)", response.BonusType)
	assert.Equal(t, "Upcoming", response.Status)
}

func TestToTournamentResponse_BonusMultiplierCalculation(t *testing.T) {
	// Given: A tournament with known base and calculated prize pool
	tournament := &models.Tournament{
		Model:               gorm.Model{ID: 1},
		Name:                "Bonus Test",
		BasePrizePool:       1000.00,
		CalculatedPrizePool: 2200.00,
		Game:                models.Game{Name: "Test Game"},
	}

	// When: Converting to response
	response := ToTournamentResponse(tournament)

	// Then: Bonus multiplier should be calculated correctly (2.2x)
	assert.Equal(t, 2.2, response.BonusMultiplier)
}

func TestToTournamentResponse_ZeroBasePrizePool(t *testing.T) {
	// Given: A tournament with zero base prize pool
	tournament := &models.Tournament{
		Model:               gorm.Model{ID: 1},
		Name:                "Free Tournament",
		BasePrizePool:       0.00,
		CalculatedPrizePool: 0.00,
		Game:                models.Game{Name: "Free Game"},
	}

	// When: Converting to response
	response := ToTournamentResponse(tournament)

	// Then: Bonus multiplier should default to 1.0
	assert.Equal(t, 1.0, response.BonusMultiplier)
}

func TestToTournamentModel_BasicRequest(t *testing.T) {
	// Given: A create tournament request
	startDate := time.Date(2024, 8, 1, 14, 0, 0, 0, time.UTC)
	req := dtos.CreateTournamentRequest{
		Name:      "New Tournament",
		GameId:    10,
		PrizePool: 5000.00,
		StartDate: startDate,
	}

	// When: Converting to model
	tournament := ToTournamentModel(req)

	// Then: Fields should be set correctly
	assert.Equal(t, "New Tournament", tournament.Name)
	assert.Equal(t, uint(10), tournament.GameID)
	assert.Equal(t, 5000.00, tournament.BasePrizePool)
	assert.Equal(t, startDate, tournament.StartDate)
	assert.Equal(t, models.StatusUpcoming, tournament.Status)
}

func TestToTournamentModel_AppliesStrategy(t *testing.T) {
	// Given: A create tournament request with a July date (summer strategy)
	startDate := time.Date(2024, 7, 15, 10, 0, 0, 0, time.UTC)
	req := dtos.CreateTournamentRequest{
		Name:      "Summer Tournament",
		GameId:    1,
		PrizePool: 1000.00,
		StartDate: startDate,
	}

	// When: Converting to model
	tournament := ToTournamentModel(req)

	// Then: The calculated prize pool should have summer bonus applied (1.2x)
	assert.Equal(t, 1200.00, tournament.CalculatedPrizePool)
	assert.Equal(t, "Summer Bonus (20%)", tournament.BonusType)
}

func TestToTournamentModel_ChristmasStrategy(t *testing.T) {
	// Given: A create tournament request with a December 25 date
	startDate := time.Date(2024, 12, 25, 12, 0, 0, 0, time.UTC)
	req := dtos.CreateTournamentRequest{
		Name:      "Christmas Tournament",
		GameId:    1,
		PrizePool: 1000.00,
		StartDate: startDate,
	}

	// When: Converting to model
	tournament := ToTournamentModel(req)

	// Then: The calculated prize pool should have Christmas bonus applied (2.2x)
	assert.Equal(t, 2200.00, tournament.CalculatedPrizePool)
	assert.Equal(t, "Christmas Bonus (120%)", tournament.BonusType)
}

func TestToTournamentModel_NormalStrategy(t *testing.T) {
	// Given: A create tournament request with a February date (normal period)
	startDate := time.Date(2024, 2, 15, 10, 0, 0, 0, time.UTC)
	req := dtos.CreateTournamentRequest{
		Name:      "Normal Tournament",
		GameId:    1,
		PrizePool: 1000.00,
		StartDate: startDate,
	}

	// When: Converting to model
	tournament := ToTournamentModel(req)

	// Then: The calculated prize pool should be the same as base (1.0x)
	assert.Equal(t, 1000.00, tournament.CalculatedPrizePool)
	assert.Equal(t, "Normal", tournament.BonusType)
}

func TestToTournamentResponseList_EmptySlice(t *testing.T) {
	// Given: An empty slice of tournaments
	tournaments := []models.Tournament{}

	// When: Converting to response list
	responses := ToTournamentResponseList(tournaments)

	// Then: The result should be an empty slice
	assert.Empty(t, responses)
	assert.Len(t, responses, 0)
}

func TestToTournamentResponseList_MultipleTournaments(t *testing.T) {
	// Given: A slice with multiple tournaments
	tournaments := []models.Tournament{
		{Model: gorm.Model{ID: 1}, Name: "Tournament 1", Game: models.Game{Name: "Game 1"}, BasePrizePool: 1000},
		{Model: gorm.Model{ID: 2}, Name: "Tournament 2", Game: models.Game{Name: "Game 2"}, BasePrizePool: 2000},
		{Model: gorm.Model{ID: 3}, Name: "Tournament 3", Game: models.Game{Name: "Game 3"}, BasePrizePool: 3000},
	}

	// When: Converting to response list
	responses := ToTournamentResponseList(tournaments)

	// Then: All tournaments should be converted
	assert.Len(t, responses, 3)
	assert.Equal(t, "Tournament 1", responses[0].Name)
	assert.Equal(t, "Tournament 2", responses[1].Name)
	assert.Equal(t, "Tournament 3", responses[2].Name)
}

func TestUpdateTournamentFromRequest_UpdatesName(t *testing.T) {
	// Given: An existing tournament and an update request
	existingTournament := &models.Tournament{
		Model:     gorm.Model{ID: 10},
		Name:      "Old Name",
		GameID:    1,
		StartDate: time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC),
	}
	req := dtos.CreateTournamentRequest{
		Name:      "New Name",
		GameId:    2,
		PrizePool: 5000.00,
		StartDate: time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC),
	}

	// When: Updating the tournament
	updatedTournament := UpdateTournamentFromRequest(existingTournament, req)

	// Then: The name should be updated
	assert.Equal(t, "New Name", updatedTournament.Name)
	assert.Equal(t, uint(2), updatedTournament.GameID)
}

func TestUpdateTournamentFromRequest_ReappliesStrategyOnDateChange(t *testing.T) {
	// Given: An existing tournament with normal date and a request with July date
	existingTournament := &models.Tournament{
		Model:               gorm.Model{ID: 10},
		Name:                "Tournament",
		BasePrizePool:       1000.00,
		CalculatedPrizePool: 1000.00,
		BonusType:           "Normal",
		StartDate:           time.Date(2024, 3, 1, 10, 0, 0, 0, time.UTC),
	}
	req := dtos.CreateTournamentRequest{
		Name:      "Tournament",
		GameId:    1,
		PrizePool: 1000.00,
		StartDate: time.Date(2024, 7, 15, 10, 0, 0, 0, time.UTC),
	}

	// When: Updating the tournament with new date
	updatedTournament := UpdateTournamentFromRequest(existingTournament, req)

	// Then: The strategy should be reapplied with summer bonus
	assert.Equal(t, 1200.00, updatedTournament.CalculatedPrizePool)
	assert.Equal(t, "Summer Bonus (20%)", updatedTournament.BonusType)
}

func TestUpdateTournamentFromRequest_UpdatesPrizePool(t *testing.T) {
	// Given: An existing tournament and a request with different prize pool
	existingTournament := &models.Tournament{
		Model:               gorm.Model{ID: 10},
		Name:                "Tournament",
		BasePrizePool:       1000.00,
		CalculatedPrizePool: 1000.00,
		BonusType:           "Normal",
		StartDate:           time.Date(2024, 3, 1, 10, 0, 0, 0, time.UTC),
	}
	req := dtos.CreateTournamentRequest{
		Name:      "Tournament",
		GameId:    1,
		PrizePool: 2000.00,
		StartDate: time.Date(2024, 3, 1, 10, 0, 0, 0, time.UTC),
	}

	// When: Updating the tournament with new prize pool
	updatedTournament := UpdateTournamentFromRequest(existingTournament, req)

	// Then: The base prize pool should be updated
	assert.Equal(t, 2000.00, updatedTournament.BasePrizePool)
}

func TestToTournamentResponse_AllStatusTypes(t *testing.T) {
	// Given: A tournament with Active status
	tournament := &models.Tournament{
		Model:  gorm.Model{ID: 1},
		Name:   "Active Tournament",
		Status: models.StatusActive,
		Game:   models.Game{Name: "Game"},
	}

	// When: Converting to response
	response := ToTournamentResponse(tournament)

	// Then: Status should be converted to string "Active"
	assert.Equal(t, "Active", response.Status)
}

func TestToTournamentResponse_CompletedStatus(t *testing.T) {
	// Given: A tournament with Completed status
	tournament := &models.Tournament{
		Model:  gorm.Model{ID: 1},
		Name:   "Completed Tournament",
		Status: models.StatusCompleted,
		Game:   models.Game{Name: "Game"},
	}

	// When: Converting to response
	response := ToTournamentResponse(tournament)

	// Then: Status should be converted to string "Completed"
	assert.Equal(t, "Completed", response.Status)
}
