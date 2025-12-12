package models

import (
	"gorm.io/gorm"
)

type Game struct {
	gorm.Model
	Name            string `gorm:"not null;unique"`
	Description     string `gorm:"type:text"`
	NumberOfPlayers int    `gorm:"not null"`

	Tournaments []Tournament `gorm:"foreignKey:GameID"`
}
