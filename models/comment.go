package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Content string `gorm:"type:text;not null"`
	UserID  uint   `gorm:"not null"`
	NewsID  uint   `gorm:"not null"`

	User User `gorm:"foreignKey:UserID"`
	News News `gorm:"foreignKey:NewsID"`
}
