package models

import (
	"testing"
	"time"

	"github.com/PI-Team04-GameClub/gameclub-backend/observer"
	"github.com/stretchr/testify/assert"
)

func TestTournament_ApplyPrizePoolStrategy_NormalPeriod(t *testing.T) {
	// Given: A tournament with a start date in February (normal period)
	tournament := &Tournament{
		Name:          "Test Tournament",
		BasePrizePool: 1000.00,
		StartDate:     time.Date(2024, 2, 15, 10, 0, 0, 0, time.UTC),
	}

	// When: Applying the prize pool strategy
	tournament.ApplyPrizePoolStrategy()

	// Then: The calculated prize pool should be the same as base (1.0x)
	assert.Equal(t, 1000.00, tournament.CalculatedPrizePool)
	assert.Equal(t, "Normal", tournament.BonusType)
}

func TestTournament_ApplyPrizePoolStrategy_SummerPeriod(t *testing.T) {
	// Given: A tournament with a start date in July (summer period)
	tournament := &Tournament{
		Name:          "Summer Tournament",
		BasePrizePool: 1000.00,
		StartDate:     time.Date(2024, 7, 15, 10, 0, 0, 0, time.UTC),
	}

	// When: Applying the prize pool strategy
	tournament.ApplyPrizePoolStrategy()

	// Then: The calculated prize pool should be 1.2x the base
	assert.Equal(t, 1200.00, tournament.CalculatedPrizePool)
	assert.Equal(t, "Summer Bonus (20%)", tournament.BonusType)
}

func TestTournament_ApplyPrizePoolStrategy_ChristmasPeriod(t *testing.T) {
	// Given: A tournament with a start date on December 25 (Christmas period)
	tournament := &Tournament{
		Name:          "Christmas Tournament",
		BasePrizePool: 1000.00,
		StartDate:     time.Date(2024, 12, 25, 10, 0, 0, 0, time.UTC),
	}

	// When: Applying the prize pool strategy
	tournament.ApplyPrizePoolStrategy()

	// Then: The calculated prize pool should be 2.2x the base
	assert.Equal(t, 2200.00, tournament.CalculatedPrizePool)
	assert.Equal(t, "Christmas Bonus (120%)", tournament.BonusType)
}

func TestTournament_GetPrizePoolBonus_ZeroBase(t *testing.T) {
	// Given: A tournament with zero base prize pool
	tournament := &Tournament{
		BasePrizePool:       0.00,
		CalculatedPrizePool: 0.00,
	}

	// When: Getting the prize pool bonus
	bonus := tournament.GetPrizePoolBonus()

	// Then: The bonus should default to 1.0
	assert.Equal(t, 1.0, bonus)
}

func TestTournament_GetPrizePoolBonus_NormalBonus(t *testing.T) {
	// Given: A tournament with equal base and calculated prize pool
	tournament := &Tournament{
		BasePrizePool:       1000.00,
		CalculatedPrizePool: 1000.00,
	}

	// When: Getting the prize pool bonus
	bonus := tournament.GetPrizePoolBonus()

	// Then: The bonus should be 1.0
	assert.Equal(t, 1.0, bonus)
}

func TestTournament_GetPrizePoolBonus_SummerBonus(t *testing.T) {
	// Given: A tournament with summer bonus applied
	tournament := &Tournament{
		BasePrizePool:       1000.00,
		CalculatedPrizePool: 1200.00,
	}

	// When: Getting the prize pool bonus
	bonus := tournament.GetPrizePoolBonus()

	// Then: The bonus should be 1.2
	assert.Equal(t, 1.2, bonus)
}

func TestTournament_GetPrizePoolBonus_ChristmasBonus(t *testing.T) {
	// Given: A tournament with Christmas bonus applied
	tournament := &Tournament{
		BasePrizePool:       1000.00,
		CalculatedPrizePool: 2200.00,
	}

	// When: Getting the prize pool bonus
	bonus := tournament.GetPrizePoolBonus()

	// Then: The bonus should be 2.2
	assert.Equal(t, 2.2, bonus)
}

type mockObserver struct {
	calledWith observer.TournamentData
	callCount  int
}

func (m *mockObserver) OnTournamentCreated(data observer.TournamentData) {
	m.calledWith = data
	m.callCount++
}

func TestTournament_Attach(t *testing.T) {
	// Given: A tournament and an observer
	tournament := &Tournament{Name: "Test"}
	obs := &mockObserver{}

	// When: Attaching the observer
	tournament.Attach(obs)

	// Then: The observer should be in the observers list
	assert.Len(t, tournament.observers, 1)
}

func TestTournament_Attach_MultipleObservers(t *testing.T) {
	// Given: A tournament and multiple observers
	tournament := &Tournament{Name: "Test"}
	obs1 := &mockObserver{}
	obs2 := &mockObserver{}

	// When: Attaching multiple observers
	tournament.Attach(obs1)
	tournament.Attach(obs2)

	// Then: Both observers should be in the list
	assert.Len(t, tournament.observers, 2)
}

func TestTournament_Detach(t *testing.T) {
	// Given: A tournament with an attached observer
	tournament := &Tournament{Name: "Test"}
	obs := &mockObserver{}
	tournament.Attach(obs)

	// When: Detaching the observer
	tournament.Detach(obs)

	// Then: The observer should be removed from the list
	assert.Len(t, tournament.observers, 0)
}

func TestTournament_Detach_NonexistentObserver(t *testing.T) {
	// Given: A tournament with an attached observer
	tournament := &Tournament{Name: "Test"}
	obs1 := &mockObserver{}
	obs2 := &mockObserver{}
	tournament.Attach(obs1)

	// When: Detaching an observer that was never attached
	tournament.Detach(obs2)

	// Then: The original observer should still be in the list
	assert.Len(t, tournament.observers, 1)
}

func TestTournament_NotifyCreated(t *testing.T) {
	// Given: A tournament with an attached observer
	tournament := &Tournament{
		Name:                "Notification Test",
		StartDate:           time.Date(2024, 8, 15, 10, 0, 0, 0, time.UTC),
		CalculatedPrizePool: 5000.00,
	}
	obs := &mockObserver{}
	tournament.Attach(obs)

	// When: Notifying about creation
	tournament.NotifyCreated()

	// Then: The observer should be called with correct data
	assert.Equal(t, 1, obs.callCount)
	assert.Equal(t, "Notification Test", obs.calledWith.Name)
	assert.Equal(t, 5000.00, obs.calledWith.PrizePool)
}

func TestTournament_NotifyCreated_MultipleObservers(t *testing.T) {
	// Given: A tournament with multiple observers
	tournament := &Tournament{
		Name:                "Multi Observer Test",
		StartDate:           time.Date(2024, 8, 15, 10, 0, 0, 0, time.UTC),
		CalculatedPrizePool: 3000.00,
	}
	obs1 := &mockObserver{}
	obs2 := &mockObserver{}
	tournament.Attach(obs1)
	tournament.Attach(obs2)

	// When: Notifying about creation
	tournament.NotifyCreated()

	// Then: Both observers should be called
	assert.Equal(t, 1, obs1.callCount)
	assert.Equal(t, 1, obs2.callCount)
}

func TestTournament_NotifyCreated_NoObservers(t *testing.T) {
	// Given: A tournament with no observers
	tournament := &Tournament{
		Name:                "No Observer Test",
		StartDate:           time.Date(2024, 8, 15, 10, 0, 0, 0, time.UTC),
		CalculatedPrizePool: 1000.00,
	}

	// When: Notifying about creation (should not panic)
	tournament.NotifyCreated()

	// Then: No panic should occur
	assert.Empty(t, tournament.observers)
}

func TestTournamentStatus_Values(t *testing.T) {
	// Given: Tournament status constants

	// When: Checking their values

	// Then: They should have the correct string values
	assert.Equal(t, TournamentStatus("Active"), StatusActive)
	assert.Equal(t, TournamentStatus("Upcoming"), StatusUpcoming)
	assert.Equal(t, TournamentStatus("Completed"), StatusCompleted)
}
