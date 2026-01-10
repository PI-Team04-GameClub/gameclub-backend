package observer

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEmailNotifier(t *testing.T) {
	// Given: A map of user emails
	userEmails := map[string]string{
		"user@example.com": "John",
	}

	// When: Creating a new email notifier
	notifier := NewEmailNotifier(userEmails)

	// Then: The notifier should not be nil
	assert.NotNil(t, notifier)
}

func TestNewEmailNotifier_EmptyMap(t *testing.T) {
	// Given: An empty map of user emails
	userEmails := map[string]string{}

	// When: Creating a new email notifier
	notifier := NewEmailNotifier(userEmails)

	// Then: The notifier should not be nil
	assert.NotNil(t, notifier)
}

func TestEmailNotifier_OnTournamentCreated_SingleUser(t *testing.T) {
	// Given: An email notifier with one user and tournament data
	userEmails := map[string]string{
		"john@example.com": "John",
	}
	notifier := NewEmailNotifier(userEmails)
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

	// Then: The log should contain the email being sent to the user
	output := buf.String()
	assert.Contains(t, output, "john@example.com")
	assert.Contains(t, output, "New Tournament Created!")
	assert.Contains(t, output, "1 users")
}

func TestEmailNotifier_OnTournamentCreated_MultipleUsers(t *testing.T) {
	// Given: An email notifier with multiple users
	userEmails := map[string]string{
		"john@example.com": "John",
		"jane@example.com": "Jane",
		"bob@example.com":  "Bob",
	}
	notifier := NewEmailNotifier(userEmails)
	tournamentData := TournamentData{
		Name:      "Winter Cup",
		StartDate: "2024-12-25 12:00",
		PrizePool: 10000.00,
	}

	// Capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	// When: The tournament created notification is triggered
	notifier.OnTournamentCreated(tournamentData)

	// Then: The log should indicate emails were sent to 3 users
	output := buf.String()
	assert.Contains(t, output, "3 users")
}

func TestEmailNotifier_OnTournamentCreated_EmptyUsers(t *testing.T) {
	// Given: An email notifier with no users
	userEmails := map[string]string{}
	notifier := NewEmailNotifier(userEmails)
	tournamentData := TournamentData{
		Name:      "Empty Tournament",
		StartDate: "2024-06-01 09:00",
		PrizePool: 1000.00,
	}

	// Capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	// When: The tournament created notification is triggered
	notifier.OnTournamentCreated(tournamentData)

	// Then: The log should indicate 0 users received emails
	output := buf.String()
	assert.Contains(t, output, "0 users")
}

func TestEmailNotifier_ImplementsTournamentObserver(t *testing.T) {
	// Given: An email notifier
	userEmails := map[string]string{"test@test.com": "Test"}

	// When: Checking if it implements TournamentObserver interface
	var observer TournamentObserver = NewEmailNotifier(userEmails)

	// Then: The assignment should succeed and observer should not be nil
	assert.NotNil(t, observer)
}

func TestEmailNotifier_FormatEmail(t *testing.T) {
	// Given: An email notifier and tournament data
	userEmails := map[string]string{
		"user@example.com": "TestUser",
	}
	notifier := NewEmailNotifier(userEmails)
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

	// Then: The log should contain personalized email content
	output := buf.String()
	assert.Contains(t, output, "TestUser")
	assert.Contains(t, output, "Test Tournament")
	assert.Contains(t, output, "3000.00")
}

func TestEmailNotifier_OnTournamentCreated_EmailContainsAllDetails(t *testing.T) {
	// Given: An email notifier with a user
	userEmails := map[string]string{
		"alice@example.com": "Alice",
	}
	notifier := NewEmailNotifier(userEmails)
	tournamentData := TournamentData{
		Name:      "Grand Finals",
		StartDate: "2024-11-30 18:00",
		PrizePool: 25000.00,
	}

	// Capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	// When: The tournament created notification is triggered
	notifier.OnTournamentCreated(tournamentData)

	// Then: The email body should contain tournament name, prize pool, and start date
	output := buf.String()
	assert.True(t, strings.Contains(output, "Grand Finals"))
	assert.True(t, strings.Contains(output, "25000.00"))
	assert.True(t, strings.Contains(output, "2024-11-30 18:00"))
	assert.True(t, strings.Contains(output, "GameClub Team"))
}

func TestEmailNotifier_OnTournamentCreated_LogsEmailSuccess(t *testing.T) {
	// Given: An email notifier with a user
	userEmails := map[string]string{
		"test@test.com": "Tester",
	}
	notifier := NewEmailNotifier(userEmails)
	tournamentData := TournamentData{
		Name:      "Quick Match",
		StartDate: "2024-03-15 10:00",
		PrizePool: 500.00,
	}

	// Capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	// When: The tournament created notification is triggered
	notifier.OnTournamentCreated(tournamentData)

	// Then: The log should indicate email was sent successfully
	output := buf.String()
	assert.Contains(t, output, "Email sent successfully")
}
