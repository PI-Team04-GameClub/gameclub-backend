package models

import "gorm.io/gorm"

type FriendRequestStatus string

const (
	StatusPending  FriendRequestStatus = "Pending"
	StatusAccepted FriendRequestStatus = "Accepted"
	StatusDeclined FriendRequestStatus = "Declined"
)

type FriendRequest struct {
	gorm.Model
	SenderID   uint                `gorm:"not null"`
	ReceiverID uint                `gorm:"not null"`
	Status     FriendRequestStatus `gorm:"type:varchar(20);default:'Pending'"`

	Sender   User `gorm:"foreignKey:SenderID"`
	Receiver User `gorm:"foreignKey:ReceiverID"`
}
