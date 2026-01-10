package observer

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogNotifier(t *testing.T) {
	// Given: Nothing (creating a new log notifier)

	// When: Creating a new log notifier
	notifier := NewLogNotifier()

	// Then: The notifier should not be nil
	assert.NotNil(t, notifier)
}

func TestLogNotifier_OnTournamentCreated(t *testing.T) {
	// Given: A log notifier and tournament data
	notifier := NewLogNotifier()
	tournamentData := TournamentData{
		Name:      "Summer Championship",
		StartDate: "2024-07-15 10:00",
		PrizePool: 5000.00,
	}

	// Capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	// When: The tournament created notification is triggered
	notifier.OnTournamentCreated(tournamentData)

	// Then: The log should contain the tournament information
	output := buf.String()
	assert.Contains(t, output, "TOURNAMENT CREATED")
	assert.Contains(t, output, "Summer Championship")
	assert.Contains(t, output, "5000.00")
	assert.Contains(t, output, "2024-07-15 10:00")
}

func TestLogNotifier_OnTournamentCreated_ZeroPrizePool(t *testing.T) {
	// Given: A log notifier and tournament data with zero prize pool
	notifier := NewLogNotifier()
	tournamentData := TournamentData{
		Name:      "Free Tournament",
		StartDate: "2024-08-01 14:00",
		PrizePool: 0.0,
	}

	// Capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	// When: The tournament created notification is triggered
	notifier.OnTournamentCreated(tournamentData)

	// Then: The log should contain the tournament information with zero prize pool
	output := buf.String()
	assert.Contains(t, output, "Free Tournament")
	assert.Contains(t, output, "0.00")
}

func TestLogNotifier_OnTournamentCreated_LargePrizePool(t *testing.T) {
	// Given: A log notifier and tournament data with large prize pool
	notifier := NewLogNotifier()
	tournamentData := TournamentData{
		Name:      "Grand Championship",
		StartDate: "2024-12-25 12:00",
		PrizePool: 1000000.50,
	}

	// Capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	// When: The tournament created notification is triggered
	notifier.OnTournamentCreated(tournamentData)

	// Then: The log should contain the tournament information with large prize pool
	output := buf.String()
	assert.Contains(t, output, "Grand Championship")
	assert.Contains(t, output, "1000000.50")
}

func TestLogNotifier_OnTournamentCreated_SpecialCharactersInName(t *testing.T) {
	// Given: A log notifier and tournament data with special characters in name
	notifier := NewLogNotifier()
	tournamentData := TournamentData{
		Name:      "Tournament #1 - 2024 (Special Edition)",
		StartDate: "2024-06-15 09:00",
		PrizePool: 2500.00,
	}

	// Capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	// When: The tournament created notification is triggered
	notifier.OnTournamentCreated(tournamentData)

	// Then: The log should contain the full tournament name with special characters
	output := buf.String()
	assert.Contains(t, output, "Tournament #1 - 2024 (Special Edition)")
}

func TestLogNotifier_ImplementsTournamentObserver(t *testing.T) {
	// Given: A log notifier

	// When: Checking if it implements TournamentObserver interface
	var observer TournamentObserver = NewLogNotifier()

	// Then: The assignment should succeed (no compilation error) and observer should not be nil
	assert.NotNil(t, observer)
}

func TestLogNotifier_OnTournamentCreated_LogFormat(t *testing.T) {
	// Given: A log notifier and tournament data
	notifier := NewLogNotifier()
	tournamentData := TournamentData{
		Name:      "Test Tournament",
		StartDate: "2024-05-20 16:30",
		PrizePool: 3000.00,
	}

	// Capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	// When: The tournament created notification is triggered
	notifier.OnTournamentCreated(tournamentData)

	// Then: The log should have the correct format with Name, Prize Pool, and Start
	output := buf.String()
	assert.True(t, strings.Contains(output, "Name:"))
	assert.True(t, strings.Contains(output, "Prize Pool:"))
	assert.True(t, strings.Contains(output, "Start:"))
}
