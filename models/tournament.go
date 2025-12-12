package models

import (
	"time"

	"gorm.io/gorm"
)

type TournamentStatus string

const (
	StatusActive    TournamentStatus = "Active"
	StatusUpcoming  TournamentStatus = "Upcoming"
	StatusCompleted TournamentStatus = "Completed"
)

type Tournament struct {
	gorm.Model
	Name      string `gorm:"not null"`
	GameID    uint   `gorm:"not null"`
	PrizePool float64
	StartDate time.Time
	Status    TournamentStatus `gorm:"type:varchar(20);default:'Upcoming'"`

	Game  Game    `gorm:"foreignKey:GameID"`
	Teams []*Team `gorm:"many2many:team_tournaments;"`
}
