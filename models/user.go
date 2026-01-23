package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string `gorm:"not null"`
	LastName  string
	Email     string `gorm:"unique"`
	Password  string `gorm:"not null"`

	Teams    []*Team   `gorm:"many2many:user_teams;"`
	News     []News    `gorm:"foreignKey:AuthorID"`
	Comments []Comment `gorm:"foreignKey:UserID"`
}
