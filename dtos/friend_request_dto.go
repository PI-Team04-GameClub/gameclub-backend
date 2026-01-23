package dtos

type CreateFriendRequestRequest struct {
	SenderID   uint `json:"sender_id" validate:"required"`
	ReceiverID uint `json:"receiver_id" validate:"required"`
}

type FriendRequestResponse struct {
	ID           uint   `json:"id"`
	SenderID     uint   `json:"sender_id"`
	SenderName   string `json:"sender_name"`
	ReceiverID   uint   `json:"receiver_id"`
	ReceiverName string `json:"receiver_name"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type FriendResponse struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
