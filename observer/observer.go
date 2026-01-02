package observer

// TournamentData contains essential tournament information for notifications
type TournamentData struct {
	Name      string
	StartDate string
	PrizePool float64
}

// TournamentObserver defines the interface for objects that want to be notified of tournament events
type TournamentObserver interface {
	OnTournamentCreated(tournament TournamentData)
}
