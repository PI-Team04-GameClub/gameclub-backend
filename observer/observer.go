package observer

type TournamentData struct {
	Name      string
	StartDate string
	PrizePool float64
}

type TournamentObserver interface {
	OnTournamentCreated(tournament TournamentData)
}
