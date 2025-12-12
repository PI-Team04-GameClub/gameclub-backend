package models

import (
	"gorm.io/gorm"
)

type Team struct {
	gorm.Model
	Name string `gorm:"not null"`

	Users       []*User       `gorm:"many2many:user_teams;"`
	Tournaments []*Tournament `gorm:"many2many:team_tournaments;"`
}
