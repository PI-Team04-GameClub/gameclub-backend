package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupFriendRequestTestApp(db *gorm.DB) *fiber.App {
	app := fiber.New()
	friendRequestHandler := NewFriendRequestHandler(db)

	app.Post("/friend-requests", friendRequestHandler.CreateFriendRequest)
	app.Put("/friend-requests/:id/accept", friendRequestHandler.AcceptFriendRequest)
	app.Put("/friend-requests/:id/decline", friendRequestHandler.DeclineFriendRequest)
	app.Delete("/friend-requests/:id", friendRequestHandler.DeleteFriendRequest)
	app.Get("/users/:id/friend-requests/sent", friendRequestHandler.GetSentFriendRequests)
	app.Get("/users/:id/friend-requests/received", friendRequestHandler.GetReceivedFriendRequests)
	app.Get("/users/:id/friends", friendRequestHandler.GetFriends)

	return app
}

func createTestUsers(db *gorm.DB) (models.User, models.User) {
	user1 := models.User{FirstName: "John", LastName: "Doe", Email: "john@test.com", Password: "pass"}
	db.Create(&user1)

	user2 := models.User{FirstName: "Jane", LastName: "Smith", Email: "jane@test.com", Password: "pass"}
	db.Create(&user2)

	return user1, user2
}

func TestFriendRequestHandler_CreateFriendRequest_Success(t *testing.T) {
	// Given: Two existing users
	db := setupTestDB(t)
	db.AutoMigrate(&models.FriendRequest{})
	app := setupFriendRequestTestApp(db)

	user1, user2 := createTestUsers(db)

	reqBody := dtos.CreateFriendRequestRequest{
		SenderID:   user1.ID,
		ReceiverID: user2.ID,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/friend-requests", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create friend request
	resp, err := app.Test(req)

	// Then: The friend request should be created
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var response dtos.FriendRequestResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, user1.ID, response.SenderID)
	assert.Equal(t, user2.ID, response.ReceiverID)
	assert.Equal(t, "Pending", response.Status)
	assert.Equal(t, "John Doe", response.SenderName)
	assert.Equal(t, "Jane Smith", response.ReceiverName)
}

func TestFriendRequestHandler_CreateFriendRequest_InvalidJSON(t *testing.T) {
	// Given: An invalid JSON request body
	db := setupTestDB(t)
	db.AutoMigrate(&models.FriendRequest{})
	app := setupFriendRequestTestApp(db)

	req := httptest.NewRequest("POST", "/friend-requests", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create friend request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestFriendRequestHandler_CreateFriendRequest_MissingSenderID(t *testing.T) {
	// Given: A request without sender ID
	db := setupTestDB(t)
	db.AutoMigrate(&models.FriendRequest{})
	app := setupFriendRequestTestApp(db)

	_, user2 := createTestUsers(db)

	reqBody := dtos.CreateFriendRequestRequest{
		ReceiverID: user2.ID,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/friend-requests", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create friend request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestFriendRequestHandler_CreateFriendRequest_MissingReceiverID(t *testing.T) {
	// Given: A request without receiver ID
	db := setupTestDB(t)
	db.AutoMigrate(&models.FriendRequest{})
	app := setupFriendRequestTestApp(db)

	user1, _ := createTestUsers(db)

	reqBody := dtos.CreateFriendRequestRequest{
		SenderID: user1.ID,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/friend-requests", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create friend request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestFriendRequestHandler_CreateFriendRequest_SelfRequest(t *testing.T) {
	// Given: A request to send friend request to self
	db := setupTestDB(t)
	db.AutoMigrate(&models.FriendRequest{})
	app := setupFriendRequestTestApp(db)

	user1, _ := createTestUsers(db)

	reqBody := dtos.CreateFriendRequestRequest{
		SenderID:   user1.ID,
		ReceiverID: user1.ID,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/friend-requests", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create friend request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestFriendRequestHandler_CreateFriendRequest_SenderNotFound(t *testing.T) {
	// Given: A non-existent sender
	db := setupTestDB(t)
	db.AutoMigrate(&models.FriendRequest{})
	app := setupFriendRequestTestApp(db)

	_, user2 := createTestUsers(db)

	reqBody := dtos.CreateFriendRequestRequest{
		SenderID:   999,
		ReceiverID: user2.ID,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/friend-requests", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create friend request
	resp, err := app.Test(req)

	// Then: The request should fail with not found
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestFriendRequestHandler_CreateFriendRequest_ReceiverNotFound(t *testing.T) {
	// Given: A non-existent receiver
	db := setupTestDB(t)
	db.AutoMigrate(&models.FriendRequest{})
	app := setupFriendRequestTestApp(db)

	user1, _ := createTestUsers(db)

	reqBody := dtos.CreateFriendRequestRequest{
		SenderID:   user1.ID,
		ReceiverID: 999,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/friend-requests", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create friend request
	resp, err := app.Test(req)

	// Then: The request should fail with not found
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestFriendRequestHandler_CreateFriendRequest_AlreadyPending(t *testing.T) {
	// Given: A pending friend request already exists
	db := setupTestDB(t)
	db.AutoMigrate(&models.FriendRequest{})
	app := setupFriendRequestTestApp(db)

	user1, user2 := createTestUsers(db)
	db.Create(&models.FriendRequest{SenderID: user1.ID, ReceiverID: user2.ID, Status: models.StatusPending})

	reqBody := dtos.CreateFriendRequestRequest{
		SenderID:   user1.ID,
		ReceiverID: user2.ID,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/friend-requests", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create friend request
	resp, err := app.Test(req)

	// Then: The request should fail with conflict
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusConflict, resp.StatusCode)
}

func TestFriendRequestHandler_CreateFriendRequest_AlreadyFriends(t *testing.T) {
	// Given: Users are already friends
	db := setupTestDB(t)
	db.AutoMigrate(&models.FriendRequest{})
	app := setupFriendRequestTestApp(db)

	user1, user2 := createTestUsers(db)
	db.Create(&models.FriendRequest{SenderID: user1.ID, ReceiverID: user2.ID, Status: models.StatusAccepted})

	reqBody := dtos.CreateFriendRequestRequest{
		SenderID:   user1.ID,
		ReceiverID: user2.ID,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/friend-requests", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create friend request
	resp, err := app.Test(req)

	// Then: The request should fail with conflict
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusConflict, resp.StatusCode)
}

func TestFriendRequestHandler_CreateFriendRequest_ReverseAlreadyPending(t *testing.T) {
	// Given: A pending friend request already exists in reverse direction
	db := setupTestDB(t)
	db.AutoMigrate(&models.FriendRequest{})
	app := setupFriendRequestTestApp(db)

	user1, user2 := createTestUsers(db)
	db.Create(&models.FriendRequest{SenderID: user2.ID, ReceiverID: user1.ID, Status: models.StatusPending})

	reqBody := dtos.CreateFriendRequestRequest{
		SenderID:   user1.ID,
		ReceiverID: user2.ID,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/friend-requests", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create friend request
	resp, err := app.Test(req)

	// Then: The request should fail with conflict
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusConflict, resp.StatusCode)
}

func TestFriendRequestHandler_AcceptFriendRequest_Success(t *testing.T) {
	// Given: A pending friend request
	db := setupTestDB(t)
	db.AutoMigrate(&models.FriendRequest{})
	app := setupFriendRequestTestApp(db)

	user1, user2 := createTestUsers(db)
	friendRequest := models.FriendRequest{SenderID: user1.ID, ReceiverID: user2.ID, Status: models.StatusPending}
	db.Create(&friendRequest)

	req := httptest.NewRequest("PUT", fmt.Sprintf("/friend-requests/%d/accept", friendRequest.ID), nil)

	// When: Making the accept friend request
	resp, err := app.Test(req)

	// Then: The friend request should be accepted
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response dtos.FriendRequestResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Accepted", response.Status)
}

func TestFriendRequestHandler_AcceptFriendRequest_NotFound(t *testing.T) {
	// Given: An empty database
	db := setupTestDB(t)
	db.AutoMigrate(&models.FriendRequest{})
	app := setupFriendRequestTestApp(db)

	req := httptest.NewRequest("PUT", "/friend-requests/999/accept", nil)

	// When: Making the accept friend request
	resp, err := app.Test(req)

	// Then: The request should fail with not found
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestFriendRequestHandler_AcceptFriendRequest_InvalidID(t *testing.T) {
	// Given: An invalid friend request ID
	db := setupTestDB(t)
	db.AutoMigrate(&models.FriendRequest{})
	app := setupFriendRequestTestApp(db)

	req := httptest.NewRequest("PUT", "/friend-requests/invalid/accept", nil)

	// When: Making the accept friend request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestFriendRequestHandler_AcceptFriendRequest_NotPending(t *testing.T) {
	// Given: A friend request that is already accepted
	db := setupTestDB(t)
	db.AutoMigrate(&models.FriendRequest{})
	app := setupFriendRequestTestApp(db)

	user1, user2 := createTestUsers(db)
	friendRequest := models.FriendRequest{SenderID: user1.ID, ReceiverID: user2.ID, Status: models.StatusAccepted}
	db.Create(&friendRequest)

	req := httptest.NewRequest("PUT", fmt.Sprintf("/friend-requests/%d/accept", friendRequest.ID), nil)

	// When: Making the accept friend request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestFriendRequestHandler_DeclineFriendRequest_Success(t *testing.T) {
	// Given: A pending friend request
	db := setupTestDB(t)
	db.AutoMigrate(&models.FriendRequest{})
	app := setupFriendRequestTestApp(db)

	user1, user2 := createTestUsers(db)
	friendRequest := models.FriendRequest{SenderID: user1.ID, ReceiverID: user2.ID, Status: models.StatusPending}
	db.Create(&friendRequest)

	req := httptest.NewRequest("PUT", fmt.Sprintf("/friend-requests/%d/decline", friendRequest.ID), nil)

	// When: Making the decline friend request
	resp, err := app.Test(req)

	// Then: The friend request should be declined
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response dtos.FriendRequestResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Declined", response.Status)
}

func TestFriendRequestHandler_DeclineFriendRequest_NotFound(t *testing.T) {
	// Given: An empty database
	db := setupTestDB(t)
	db.AutoMigrate(&models.FriendRequest{})
	app := setupFriendRequestTestApp(db)

	req := httptest.NewRequest("PUT", "/friend-requests/999/decline", nil)

	// When: Making the decline friend request
	resp, err := app.Test(req)

	// Then: The request should fail with not found
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestFriendRequestHandler_DeclineFriendRequest_InvalidID(t *testing.T) {
	// Given: An invalid friend request ID
	db := setupTestDB(t)
	db.AutoMigrate(&models.FriendRequest{})
	app := setupFriendRequestTestApp(db)

	req := httptest.NewRequest("PUT", "/friend-requests/invalid/decline", nil)

	// When: Making the decline friend request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestFriendRequestHandler_DeclineFriendRequest_NotPending(t *testing.T) {
	// Given: A friend request that is already declined
	db := setupTestDB(t)
	db.AutoMigrate(&models.FriendRequest{})
	app := setupFriendRequestTestApp(db)

	user1, user2 := createTestUsers(db)
	friendRequest := models.FriendRequest{SenderID: user1.ID, ReceiverID: user2.ID, Status: models.StatusDeclined}
	db.Create(&friendRequest)

	req := httptest.NewRequest("PUT", fmt.Sprintf("/friend-requests/%d/decline", friendRequest.ID), nil)

	// When: Making the decline friend request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestFriendRequestHandler_DeleteFriendRequest_Success(t *testing.T) {
	// Given: An existing friend request
	db := setupTestDB(t)
	db.AutoMigrate(&models.FriendRequest{})
	app := setupFriendRequestTestApp(db)

	user1, user2 := createTestUsers(db)
	friendRequest := models.FriendRequest{SenderID: user1.ID, ReceiverID: user2.ID, Status: models.StatusPending}
	db.Create(&friendRequest)

	req := httptest.NewRequest("DELETE", fmt.Sprintf("/friend-requests/%d", friendRequest.ID), nil)

	// When: Making the delete friend request
	resp, err := app.Test(req)

	// Then: The friend request should be deleted
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)

	// Verify it's deleted
	var count int64
	db.Model(&models.FriendRequest{}).Where("id = ?", friendRequest.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestFriendRequestHandler_DeleteFriendRequest_NotFound(t *testing.T) {
	// Given: An empty database
	db := setupTestDB(t)
	db.AutoMigrate(&models.FriendRequest{})
	app := setupFriendRequestTestApp(db)

	req := httptest.NewRequest("DELETE", "/friend-requests/999", nil)

	// When: Making the delete friend request
	resp, err := app.Test(req)

	// Then: The request should fail with not found
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestFriendRequestHandler_DeleteFriendRequest_InvalidID(t *testing.T) {
	// Given: An invalid friend request ID
	db := setupTestDB(t)
	db.AutoMigrate(&models.FriendRequest{})
	app := setupFriendRequestTestApp(db)

	req := httptest.NewRequest("DELETE", "/friend-requests/invalid", nil)

	// When: Making the delete friend request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestFriendRequestHandler_GetSentFriendRequests_Success(t *testing.T) {
	// Given: A user with sent friend requests
	db := setupTestDB(t)
	db.AutoMigrate(&models.FriendRequest{})
	app := setupFriendRequestTestApp(db)

	user1, user2 := createTestUsers(db)
	user3 := models.User{FirstName: "Bob", LastName: "Wilson", Email: "bob@test.com", Password: "pass"}
	db.Create(&user3)

	db.Create(&models.FriendRequest{SenderID: user1.ID, ReceiverID: user2.ID, Status: models.StatusPending})
	db.Create(&models.FriendRequest{SenderID: user1.ID, ReceiverID: user3.ID, Status: models.StatusAccepted})

	req := httptest.NewRequest("GET", fmt.Sprintf("/users/%d/friend-requests/sent", user1.ID), nil)

	// When: Making the get sent friend requests request
	resp, err := app.Test(req)

	// Then: The request should return all sent friend requests
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var friendRequests []dtos.FriendRequestResponse
	json.NewDecoder(resp.Body).Decode(&friendRequests)
	assert.Len(t, friendRequests, 2)
}

func TestFriendRequestHandler_GetSentFriendRequests_Empty(t *testing.T) {
	// Given: A user with no sent friend requests
	db := setupTestDB(t)
	db.AutoMigrate(&models.FriendRequest{})
	app := setupFriendRequestTestApp(db)

	user1, _ := createTestUsers(db)

	req := httptest.NewRequest("GET", fmt.Sprintf("/users/%d/friend-requests/sent", user1.ID), nil)

	// When: Making the get sent friend requests request
	resp, err := app.Test(req)

	// Then: The request should return empty array
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var friendRequests []dtos.FriendRequestResponse
	json.NewDecoder(resp.Body).Decode(&friendRequests)
	assert.Empty(t, friendRequests)
}

func TestFriendRequestHandler_GetSentFriendRequests_UserNotFound(t *testing.T) {
	// Given: A non-existent user
	db := setupTestDB(t)
	db.AutoMigrate(&models.FriendRequest{})
	app := setupFriendRequestTestApp(db)

	req := httptest.NewRequest("GET", "/users/999/friend-requests/sent", nil)

	// When: Making the get sent friend requests request
	resp, err := app.Test(req)

	// Then: The request should fail with not found
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestFriendRequestHandler_GetSentFriendRequests_InvalidID(t *testing.T) {
	// Given: An invalid user ID
	db := setupTestDB(t)
	db.AutoMigrate(&models.FriendRequest{})
	app := setupFriendRequestTestApp(db)

	req := httptest.NewRequest("GET", "/users/invalid/friend-requests/sent", nil)

	// When: Making the get sent friend requests request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestFriendRequestHandler_GetReceivedFriendRequests_Success(t *testing.T) {
	// Given: A user with received friend requests
	db := setupTestDB(t)
	db.AutoMigrate(&models.FriendRequest{})
	app := setupFriendRequestTestApp(db)

	user1, user2 := createTestUsers(db)
	user3 := models.User{FirstName: "Bob", LastName: "Wilson", Email: "bob@test.com", Password: "pass"}
	db.Create(&user3)

	db.Create(&models.FriendRequest{SenderID: user2.ID, ReceiverID: user1.ID, Status: models.StatusPending})
	db.Create(&models.FriendRequest{SenderID: user3.ID, ReceiverID: user1.ID, Status: models.StatusPending})

	req := httptest.NewRequest("GET", fmt.Sprintf("/users/%d/friend-requests/received", user1.ID), nil)

	// When: Making the get received friend requests request
	resp, err := app.Test(req)

	// Then: The request should return only pending friend requests
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var friendRequests []dtos.FriendRequestResponse
	json.NewDecoder(resp.Body).Decode(&friendRequests)
	assert.Len(t, friendRequests, 2)
	for _, fr := range friendRequests {
		assert.Equal(t, "Pending", fr.Status)
	}
}

func TestFriendRequestHandler_GetReceivedFriendRequests_ExcludesAccepted(t *testing.T) {
	// Given: A user with mixed friend requests
	db := setupTestDB(t)
	db.AutoMigrate(&models.FriendRequest{})
	app := setupFriendRequestTestApp(db)

	user1, user2 := createTestUsers(db)
	user3 := models.User{FirstName: "Bob", LastName: "Wilson", Email: "bob@test.com", Password: "pass"}
	db.Create(&user3)

	db.Create(&models.FriendRequest{SenderID: user2.ID, ReceiverID: user1.ID, Status: models.StatusPending})
	db.Create(&models.FriendRequest{SenderID: user3.ID, ReceiverID: user1.ID, Status: models.StatusAccepted})

	req := httptest.NewRequest("GET", fmt.Sprintf("/users/%d/friend-requests/received", user1.ID), nil)

	// When: Making the get received friend requests request
	resp, err := app.Test(req)

	// Then: The request should return only pending friend requests
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var friendRequests []dtos.FriendRequestResponse
	json.NewDecoder(resp.Body).Decode(&friendRequests)
	assert.Len(t, friendRequests, 1)
	assert.Equal(t, "Pending", friendRequests[0].Status)
}

func TestFriendRequestHandler_GetReceivedFriendRequests_UserNotFound(t *testing.T) {
	// Given: A non-existent user
	db := setupTestDB(t)
	db.AutoMigrate(&models.FriendRequest{})
	app := setupFriendRequestTestApp(db)

	req := httptest.NewRequest("GET", "/users/999/friend-requests/received", nil)

	// When: Making the get received friend requests request
	resp, err := app.Test(req)

	// Then: The request should fail with not found
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestFriendRequestHandler_GetReceivedFriendRequests_InvalidID(t *testing.T) {
	// Given: An invalid user ID
	db := setupTestDB(t)
	db.AutoMigrate(&models.FriendRequest{})
	app := setupFriendRequestTestApp(db)

	req := httptest.NewRequest("GET", "/users/invalid/friend-requests/received", nil)

	// When: Making the get received friend requests request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestFriendRequestHandler_GetFriends_Success(t *testing.T) {
	// Given: A user with friends
	db := setupTestDB(t)
	db.AutoMigrate(&models.FriendRequest{})
	app := setupFriendRequestTestApp(db)

	user1, user2 := createTestUsers(db)
	user3 := models.User{FirstName: "Bob", LastName: "Wilson", Email: "bob@test.com", Password: "pass"}
	db.Create(&user3)

	// user1 sent request to user2, accepted
	db.Create(&models.FriendRequest{SenderID: user1.ID, ReceiverID: user2.ID, Status: models.StatusAccepted})
	// user3 sent request to user1, accepted
	db.Create(&models.FriendRequest{SenderID: user3.ID, ReceiverID: user1.ID, Status: models.StatusAccepted})

	req := httptest.NewRequest("GET", fmt.Sprintf("/users/%d/friends", user1.ID), nil)

	// When: Making the get friends request
	resp, err := app.Test(req)

	// Then: The request should return all friends
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var friends []dtos.FriendResponse
	json.NewDecoder(resp.Body).Decode(&friends)
	assert.Len(t, friends, 2)
}

func TestFriendRequestHandler_GetFriends_ExcludesPending(t *testing.T) {
	// Given: A user with mixed friend requests
	db := setupTestDB(t)
	db.AutoMigrate(&models.FriendRequest{})
	app := setupFriendRequestTestApp(db)

	user1, user2 := createTestUsers(db)
	user3 := models.User{FirstName: "Bob", LastName: "Wilson", Email: "bob@test.com", Password: "pass"}
	db.Create(&user3)

	db.Create(&models.FriendRequest{SenderID: user1.ID, ReceiverID: user2.ID, Status: models.StatusAccepted})
	db.Create(&models.FriendRequest{SenderID: user3.ID, ReceiverID: user1.ID, Status: models.StatusPending})

	req := httptest.NewRequest("GET", fmt.Sprintf("/users/%d/friends", user1.ID), nil)

	// When: Making the get friends request
	resp, err := app.Test(req)

	// Then: The request should return only accepted friends
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var friends []dtos.FriendResponse
	json.NewDecoder(resp.Body).Decode(&friends)
	assert.Len(t, friends, 1)
	assert.Equal(t, "Jane", friends[0].FirstName)
}

func TestFriendRequestHandler_GetFriends_Empty(t *testing.T) {
	// Given: A user with no friends
	db := setupTestDB(t)
	db.AutoMigrate(&models.FriendRequest{})
	app := setupFriendRequestTestApp(db)

	user1, _ := createTestUsers(db)

	req := httptest.NewRequest("GET", fmt.Sprintf("/users/%d/friends", user1.ID), nil)

	// When: Making the get friends request
	resp, err := app.Test(req)

	// Then: The request should return empty array
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var friends []dtos.FriendResponse
	json.NewDecoder(resp.Body).Decode(&friends)
	assert.Empty(t, friends)
}

func TestFriendRequestHandler_GetFriends_UserNotFound(t *testing.T) {
	// Given: A non-existent user
	db := setupTestDB(t)
	db.AutoMigrate(&models.FriendRequest{})
	app := setupFriendRequestTestApp(db)

	req := httptest.NewRequest("GET", "/users/999/friends", nil)

	// When: Making the get friends request
	resp, err := app.Test(req)

	// Then: The request should fail with not found
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestFriendRequestHandler_GetFriends_InvalidID(t *testing.T) {
	// Given: An invalid user ID
	db := setupTestDB(t)
	db.AutoMigrate(&models.FriendRequest{})
	app := setupFriendRequestTestApp(db)

	req := httptest.NewRequest("GET", "/users/invalid/friends", nil)

	// When: Making the get friends request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestNewFriendRequestHandler(t *testing.T) {
	// Given: A database connection
	db := setupTestDB(t)

	// When: Creating a new friend request handler
	handler := NewFriendRequestHandler(db)

	// Then: The handler should be created
	assert.NotNil(t, handler)
}
