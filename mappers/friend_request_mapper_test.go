package mappers

import (
	"testing"
	"time"

	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestToFriendRequestResponse_BasicFields(t *testing.T) {
	// Given: A friend request model with all fields populated
	now := time.Now()
	friendRequest := &models.FriendRequest{
		Model:      gorm.Model{ID: 1, CreatedAt: now, UpdatedAt: now},
		SenderID:   5,
		ReceiverID: 10,
		Status:     models.StatusPending,
		Sender:     models.User{FirstName: "John", LastName: "Doe"},
		Receiver:   models.User{FirstName: "Jane", LastName: "Smith"},
	}

	// When: Converting to response
	response := ToFriendRequestResponse(friendRequest)

	// Then: All fields should be correctly mapped
	assert.Equal(t, uint(1), response.ID)
	assert.Equal(t, uint(5), response.SenderID)
	assert.Equal(t, "John Doe", response.SenderName)
	assert.Equal(t, uint(10), response.ReceiverID)
	assert.Equal(t, "Jane Smith", response.ReceiverName)
	assert.Equal(t, "Pending", response.Status)
	assert.NotEmpty(t, response.CreatedAt)
	assert.NotEmpty(t, response.UpdatedAt)
}

func TestToFriendRequestResponse_AcceptedStatus(t *testing.T) {
	// Given: An accepted friend request
	now := time.Now()
	friendRequest := &models.FriendRequest{
		Model:      gorm.Model{ID: 2, CreatedAt: now, UpdatedAt: now},
		SenderID:   1,
		ReceiverID: 2,
		Status:     models.StatusAccepted,
		Sender:     models.User{FirstName: "John", LastName: "Doe"},
		Receiver:   models.User{FirstName: "Jane", LastName: "Smith"},
	}

	// When: Converting to response
	response := ToFriendRequestResponse(friendRequest)

	// Then: Status should be Accepted
	assert.Equal(t, "Accepted", response.Status)
}

func TestToFriendRequestResponse_DeclinedStatus(t *testing.T) {
	// Given: A declined friend request
	now := time.Now()
	friendRequest := &models.FriendRequest{
		Model:      gorm.Model{ID: 3, CreatedAt: now, UpdatedAt: now},
		SenderID:   1,
		ReceiverID: 2,
		Status:     models.StatusDeclined,
		Sender:     models.User{FirstName: "John", LastName: "Doe"},
		Receiver:   models.User{FirstName: "Jane", LastName: "Smith"},
	}

	// When: Converting to response
	response := ToFriendRequestResponse(friendRequest)

	// Then: Status should be Declined
	assert.Equal(t, "Declined", response.Status)
}

func TestToFriendRequestResponse_EmptyLastName(t *testing.T) {
	// Given: A friend request with users having only first names
	now := time.Now()
	friendRequest := &models.FriendRequest{
		Model:      gorm.Model{ID: 4, CreatedAt: now, UpdatedAt: now},
		SenderID:   1,
		ReceiverID: 2,
		Status:     models.StatusPending,
		Sender:     models.User{FirstName: "John", LastName: ""},
		Receiver:   models.User{FirstName: "Jane", LastName: ""},
	}

	// When: Converting to response
	response := ToFriendRequestResponse(friendRequest)

	// Then: Names should handle empty last name
	assert.Equal(t, "John ", response.SenderName)
	assert.Equal(t, "Jane ", response.ReceiverName)
}

func TestToFriendRequestResponse_DateTimeFormat(t *testing.T) {
	// Given: A friend request with specific timestamps
	createdAt := time.Date(2024, 7, 15, 10, 30, 45, 0, time.UTC)
	updatedAt := time.Date(2024, 7, 16, 14, 20, 30, 0, time.UTC)
	friendRequest := &models.FriendRequest{
		Model:      gorm.Model{ID: 1, CreatedAt: createdAt, UpdatedAt: updatedAt},
		SenderID:   1,
		ReceiverID: 2,
		Status:     models.StatusPending,
		Sender:     models.User{FirstName: "John", LastName: "Doe"},
		Receiver:   models.User{FirstName: "Jane", LastName: "Smith"},
	}

	// When: Converting to response
	response := ToFriendRequestResponse(friendRequest)

	// Then: Dates should be formatted correctly
	assert.Equal(t, "2024-07-15 10:30:45", response.CreatedAt)
	assert.Equal(t, "2024-07-16 14:20:30", response.UpdatedAt)
}

func TestToFriendRequestResponseList_EmptySlice(t *testing.T) {
	// Given: An empty slice of friend requests
	friendRequests := []models.FriendRequest{}

	// When: Converting to response list
	responses := ToFriendRequestResponseList(friendRequests)

	// Then: The result should be an empty slice
	assert.Empty(t, responses)
	assert.Len(t, responses, 0)
}

func TestToFriendRequestResponseList_SingleFriendRequest(t *testing.T) {
	// Given: A slice with one friend request
	now := time.Now()
	friendRequests := []models.FriendRequest{
		{
			Model:      gorm.Model{ID: 1, CreatedAt: now, UpdatedAt: now},
			SenderID:   1,
			ReceiverID: 2,
			Status:     models.StatusPending,
			Sender:     models.User{FirstName: "John", LastName: "Doe"},
			Receiver:   models.User{FirstName: "Jane", LastName: "Smith"},
		},
	}

	// When: Converting to response list
	responses := ToFriendRequestResponseList(friendRequests)

	// Then: One friend request should be converted
	assert.Len(t, responses, 1)
	assert.Equal(t, "John Doe", responses[0].SenderName)
	assert.Equal(t, "Pending", responses[0].Status)
}

func TestToFriendRequestResponseList_MultipleFriendRequests(t *testing.T) {
	// Given: A slice with multiple friend requests
	now := time.Now()
	friendRequests := []models.FriendRequest{
		{Model: gorm.Model{ID: 1, CreatedAt: now, UpdatedAt: now}, Status: models.StatusPending, Sender: models.User{FirstName: "User1"}, Receiver: models.User{FirstName: "User2"}},
		{Model: gorm.Model{ID: 2, CreatedAt: now, UpdatedAt: now}, Status: models.StatusAccepted, Sender: models.User{FirstName: "User3"}, Receiver: models.User{FirstName: "User4"}},
		{Model: gorm.Model{ID: 3, CreatedAt: now, UpdatedAt: now}, Status: models.StatusDeclined, Sender: models.User{FirstName: "User5"}, Receiver: models.User{FirstName: "User6"}},
	}

	// When: Converting to response list
	responses := ToFriendRequestResponseList(friendRequests)

	// Then: All friend requests should be converted
	assert.Len(t, responses, 3)
	assert.Equal(t, "Pending", responses[0].Status)
	assert.Equal(t, "Accepted", responses[1].Status)
	assert.Equal(t, "Declined", responses[2].Status)
}

func TestToFriendRequestResponseList_PreservesOrder(t *testing.T) {
	// Given: A slice with friend requests in specific order
	now := time.Now()
	friendRequests := []models.FriendRequest{
		{Model: gorm.Model{ID: 3, CreatedAt: now, UpdatedAt: now}, Status: models.StatusDeclined, Sender: models.User{FirstName: "Third"}, Receiver: models.User{FirstName: "R"}},
		{Model: gorm.Model{ID: 1, CreatedAt: now, UpdatedAt: now}, Status: models.StatusPending, Sender: models.User{FirstName: "First"}, Receiver: models.User{FirstName: "R"}},
		{Model: gorm.Model{ID: 2, CreatedAt: now, UpdatedAt: now}, Status: models.StatusAccepted, Sender: models.User{FirstName: "Second"}, Receiver: models.User{FirstName: "R"}},
	}

	// When: Converting to response list
	responses := ToFriendRequestResponseList(friendRequests)

	// Then: The order should be preserved
	assert.Equal(t, "Third ", responses[0].SenderName)
	assert.Equal(t, "First ", responses[1].SenderName)
	assert.Equal(t, "Second ", responses[2].SenderName)
}

func TestToFriendRequestModel_BasicRequest(t *testing.T) {
	// Given: A create friend request request
	req := dtos.CreateFriendRequestRequest{
		SenderID:   10,
		ReceiverID: 20,
	}

	// When: Converting to model
	friendRequest := ToFriendRequestModel(req)

	// Then: Fields should be set correctly
	assert.Equal(t, uint(10), friendRequest.SenderID)
	assert.Equal(t, uint(20), friendRequest.ReceiverID)
	assert.Equal(t, models.StatusPending, friendRequest.Status)
}

func TestToFriendRequestModel_DefaultStatus(t *testing.T) {
	// Given: A create friend request request
	req := dtos.CreateFriendRequestRequest{
		SenderID:   1,
		ReceiverID: 2,
	}

	// When: Converting to model
	friendRequest := ToFriendRequestModel(req)

	// Then: Status should default to Pending
	assert.Equal(t, models.StatusPending, friendRequest.Status)
}

func TestToFriendResponse_BasicFields(t *testing.T) {
	// Given: A user model with all fields populated
	user := &models.User{
		Model:     gorm.Model{ID: 1},
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@test.com",
	}

	// When: Converting to friend response
	response := ToFriendResponse(user)

	// Then: All fields should be correctly mapped
	assert.Equal(t, uint(1), response.ID)
	assert.Equal(t, "John", response.FirstName)
	assert.Equal(t, "Doe", response.LastName)
	assert.Equal(t, "john@test.com", response.Email)
}

func TestToFriendResponse_EmptyLastName(t *testing.T) {
	// Given: A user with no last name
	user := &models.User{
		Model:     gorm.Model{ID: 2},
		FirstName: "Jane",
		LastName:  "",
		Email:     "jane@test.com",
	}

	// When: Converting to friend response
	response := ToFriendResponse(user)

	// Then: Empty last name should be preserved
	assert.Equal(t, "Jane", response.FirstName)
	assert.Equal(t, "", response.LastName)
}

func TestToFriendResponseList_EmptySlice(t *testing.T) {
	// Given: An empty slice of users
	users := []models.User{}

	// When: Converting to friend response list
	responses := ToFriendResponseList(users)

	// Then: The result should be an empty slice
	assert.Empty(t, responses)
	assert.Len(t, responses, 0)
}

func TestToFriendResponseList_SingleUser(t *testing.T) {
	// Given: A slice with one user
	users := []models.User{
		{Model: gorm.Model{ID: 1}, FirstName: "John", LastName: "Doe", Email: "john@test.com"},
	}

	// When: Converting to friend response list
	responses := ToFriendResponseList(users)

	// Then: One user should be converted
	assert.Len(t, responses, 1)
	assert.Equal(t, "John", responses[0].FirstName)
	assert.Equal(t, "john@test.com", responses[0].Email)
}

func TestToFriendResponseList_MultipleUsers(t *testing.T) {
	// Given: A slice with multiple users
	users := []models.User{
		{Model: gorm.Model{ID: 1}, FirstName: "John", LastName: "Doe", Email: "john@test.com"},
		{Model: gorm.Model{ID: 2}, FirstName: "Jane", LastName: "Smith", Email: "jane@test.com"},
		{Model: gorm.Model{ID: 3}, FirstName: "Bob", LastName: "Wilson", Email: "bob@test.com"},
	}

	// When: Converting to friend response list
	responses := ToFriendResponseList(users)

	// Then: All users should be converted
	assert.Len(t, responses, 3)
	assert.Equal(t, "John", responses[0].FirstName)
	assert.Equal(t, "Jane", responses[1].FirstName)
	assert.Equal(t, "Bob", responses[2].FirstName)
}

func TestToFriendResponseList_PreservesOrder(t *testing.T) {
	// Given: A slice with users in specific order
	users := []models.User{
		{Model: gorm.Model{ID: 3}, FirstName: "Third", LastName: "User", Email: "third@test.com"},
		{Model: gorm.Model{ID: 1}, FirstName: "First", LastName: "User", Email: "first@test.com"},
		{Model: gorm.Model{ID: 2}, FirstName: "Second", LastName: "User", Email: "second@test.com"},
	}

	// When: Converting to friend response list
	responses := ToFriendResponseList(users)

	// Then: The order should be preserved
	assert.Equal(t, "Third", responses[0].FirstName)
	assert.Equal(t, "First", responses[1].FirstName)
	assert.Equal(t, "Second", responses[2].FirstName)
}
