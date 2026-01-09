package observer

import "log"

type LogNotifier struct{}

func NewLogNotifier() *LogNotifier {
	return &LogNotifier{}
}

func (l *LogNotifier) OnTournamentCreated(tournament TournamentData) {
	log.Printf(
		"TOURNAMENT CREATED - Name: %s, Prize Pool: $%.2f, Start: %s",
		tournament.Name,
		tournament.PrizePool,
		tournament.StartDate,
	)
}
