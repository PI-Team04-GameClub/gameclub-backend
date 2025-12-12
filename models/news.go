package models

import (
	"gorm.io/gorm"
)

type News struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string `gorm:"type:text"`
	AuthorID    uint   `gorm:"not null"`
	Date        string

	Author User `gorm:"foreignKey:AuthorID"`
}
