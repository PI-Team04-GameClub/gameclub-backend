package models

import (
	"time"

	"github.com/PI-Team04-GameClub/gameclub-backend/strategy"
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
	Name                string  `gorm:"not null"`
	GameID              uint    `gorm:"not null"`
	BasePrizePool       float64 `gorm:"type:decimal(10,2)"`                // Original prize pool amount
	CalculatedPrizePool float64 `gorm:"type:decimal(10,2)"`                // Prize pool after applying strategy
	BonusType           string  `gorm:"type:varchar(50);default:'Normal'"` // Which bonus was applied
	StartDate           time.Time
	Status              TournamentStatus `gorm:"type:varchar(20);default:'Upcoming'"`

	Game  Game    `gorm:"foreignKey:GameID"`
	Teams []*Team `gorm:"many2many:team_tournaments;"`
}

// ApplyPrizePoolStrategy calculates and sets the prize pool based on start date
// Uses the Strategy pattern to select and apply the appropriate calculation strategy
func (t *Tournament) ApplyPrizePoolStrategy() {
	selectedStrategy := strategy.GetStrategyForDate(t.StartDate)
	calculator := strategy.NewCalculator(selectedStrategy)

	t.CalculatedPrizePool = calculator.Calculate(t.BasePrizePool)
	t.BonusType = selectedStrategy.GetStrategyName()
}

// GetPrizePoolBonus returns the bonus multiplier that was applied
// Returns 1.0 if no bonus, 1.2 for summer, 2.2 for christmas
func (t *Tournament) GetPrizePoolBonus() float64 {
	if t.BasePrizePool == 0 {
		return 1.0
	}
	return t.CalculatedPrizePool / t.BasePrizePool
}
