package mappers

import (
	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
)

func ToUserResponse(user *models.User) dtos.UserResponse {
	return dtos.UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
}

func ToUserResponseList(users []models.User) []dtos.UserResponse {
	responses := make([]dtos.UserResponse, len(users))
	for i, user := range users {
		responses[i] = ToUserResponse(&user)
	}
	return responses
}

func ToUserResponseListFromPointers(users []*models.User) []dtos.UserResponse {
	responses := make([]dtos.UserResponse, len(users))
	for i, user := range users {
		responses[i] = ToUserResponse(user)
	}
	return responses
}

func UpdateUserFromRequest(existingUser *models.User, req dtos.UpdateUserRequest) *models.User {
	if req.FirstName != "" {
		existingUser.FirstName = req.FirstName
	}
	if req.LastName != "" {
		existingUser.LastName = req.LastName
	}
	if req.Email != "" {
		existingUser.Email = req.Email
	}
	if req.Password != "" {
		existingUser.Password = req.Password
	}
	return existingUser
}
