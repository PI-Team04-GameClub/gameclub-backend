package observer

import "log"

// LogNotifier logs tournament events for audit trail
type LogNotifier struct{}

// NewLogNotifier creates a new log notifier observer
func NewLogNotifier() *LogNotifier {
	return &LogNotifier{}
}

// OnTournamentCreated logs when a new tournament is created
func (l *LogNotifier) OnTournamentCreated(tournament TournamentData) {
	log.Printf(
		"TOURNAMENT CREATED - Name: %s, Prize Pool: $%.2f, Start: %s",
		tournament.Name,
		tournament.PrizePool,
		tournament.StartDate,
	)
}
