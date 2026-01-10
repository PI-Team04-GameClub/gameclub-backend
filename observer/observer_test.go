package observer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTournamentData_Fields(t *testing.T) {
	// Given: Tournament data with specific values
	data := TournamentData{
		Name:      "Test Tournament",
		StartDate: "2024-07-15 10:00",
		PrizePool: 5000.00,
	}

	// When: Accessing the fields

	// Then: The fields should have the correct values
	assert.Equal(t, "Test Tournament", data.Name)
	assert.Equal(t, "2024-07-15 10:00", data.StartDate)
	assert.Equal(t, 5000.00, data.PrizePool)
}

func TestTournamentData_EmptyValues(t *testing.T) {
	// Given: Tournament data with empty/zero values
	data := TournamentData{
		Name:      "",
		StartDate: "",
		PrizePool: 0.0,
	}

	// When: Accessing the fields

	// Then: The fields should have empty/zero values
	assert.Equal(t, "", data.Name)
	assert.Equal(t, "", data.StartDate)
	assert.Equal(t, 0.0, data.PrizePool)
}

// MockObserver is a test mock for TournamentObserver
type MockObserver struct {
	CalledWith TournamentData
	CallCount  int
}

func (m *MockObserver) OnTournamentCreated(tournament TournamentData) {
	m.CalledWith = tournament
	m.CallCount++
}

func TestMockObserver_ImplementsInterface(t *testing.T) {
	// Given: A mock observer

	// When: Checking if it implements TournamentObserver interface
	var observer TournamentObserver = &MockObserver{}

	// Then: The assignment should succeed
	assert.NotNil(t, observer)
}

func TestMockObserver_OnTournamentCreated(t *testing.T) {
	// Given: A mock observer and tournament data
	mockObserver := &MockObserver{}
	tournamentData := TournamentData{
		Name:      "Mock Tournament",
		StartDate: "2024-08-20 15:00",
		PrizePool: 7500.00,
	}

	// When: Calling OnTournamentCreated
	mockObserver.OnTournamentCreated(tournamentData)

	// Then: The mock should record the call
	assert.Equal(t, 1, mockObserver.CallCount)
	assert.Equal(t, tournamentData, mockObserver.CalledWith)
}

func TestMockObserver_MultipleCalls(t *testing.T) {
	// Given: A mock observer
	mockObserver := &MockObserver{}
	data1 := TournamentData{Name: "Tournament 1", StartDate: "2024-01-01 10:00", PrizePool: 1000.0}
	data2 := TournamentData{Name: "Tournament 2", StartDate: "2024-02-01 10:00", PrizePool: 2000.0}

	// When: Calling OnTournamentCreated multiple times
	mockObserver.OnTournamentCreated(data1)
	mockObserver.OnTournamentCreated(data2)

	// Then: The call count should be 2 and last call data should be stored
	assert.Equal(t, 2, mockObserver.CallCount)
	assert.Equal(t, data2, mockObserver.CalledWith)
}
