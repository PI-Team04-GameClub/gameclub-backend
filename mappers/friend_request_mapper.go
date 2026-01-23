package mappers

import (
	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
)

func ToFriendRequestResponse(fr *models.FriendRequest) dtos.FriendRequestResponse {
	return dtos.FriendRequestResponse{
		ID:           fr.ID,
		SenderID:     fr.SenderID,
		SenderName:   fr.Sender.FirstName + " " + fr.Sender.LastName,
		ReceiverID:   fr.ReceiverID,
		ReceiverName: fr.Receiver.FirstName + " " + fr.Receiver.LastName,
		Status:       string(fr.Status),
		CreatedAt:    fr.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    fr.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToFriendRequestResponseList(friendRequests []models.FriendRequest) []dtos.FriendRequestResponse {
	responses := make([]dtos.FriendRequestResponse, len(friendRequests))
	for i, fr := range friendRequests {
		responses[i] = ToFriendRequestResponse(&fr)
	}
	return responses
}

func ToFriendRequestModel(req dtos.CreateFriendRequestRequest) models.FriendRequest {
	return models.FriendRequest{
		SenderID:   req.SenderID,
		ReceiverID: req.ReceiverID,
		Status:     models.StatusPending,
	}
}

func ToFriendResponse(user *models.User) dtos.FriendResponse {
	return dtos.FriendResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
}

func ToFriendResponseList(users []models.User) []dtos.FriendResponse {
	responses := make([]dtos.FriendResponse, len(users))
	for i, user := range users {
		responses[i] = ToFriendResponse(&user)
	}
	return responses
}
