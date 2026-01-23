package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/mocks"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestFriendRequestHandler_CreateFriendRequest_Success_Unit(t *testing.T) {
	// Given: A valid create friend request
	mockFriendRequestRepo := new(mocks.MockFriendRequestRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewFriendRequestHandlerWithRepo(mockFriendRequestRepo, mockUserRepo)

	app := fiber.New()
	app.Post("/friend-requests", handler.CreateFriendRequest)

	sender := &models.User{Model: gorm.Model{ID: 1}, FirstName: "John", LastName: "Doe"}
	receiver := &models.User{Model: gorm.Model{ID: 2}, FirstName: "Jane", LastName: "Smith"}
	mockUserRepo.On("FindByID", mock.Anything, uint(1)).Return(sender, nil)
	mockUserRepo.On("FindByID", mock.Anything, uint(2)).Return(receiver, nil)
	mockFriendRequestRepo.On("FindByUsers", mock.Anything, uint(1), uint(2)).Return(nil, gorm.ErrRecordNotFound)
	mockFriendRequestRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.FriendRequest")).Return(nil)
	mockFriendRequestRepo.On("FindByID", mock.Anything, uint(0)).Return(&models.FriendRequest{
		Model:      gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		SenderID:   1,
		ReceiverID: 2,
		Status:     models.StatusPending,
		Sender:     *sender,
		Receiver:   *receiver,
	}, nil)

	reqBody := dtos.CreateFriendRequestRequest{
		SenderID:   1,
		ReceiverID: 2,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/friend-requests", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create friend request
	resp, err := app.Test(req)

	// Then: The request should succeed with created status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	mockFriendRequestRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestFriendRequestHandler_CreateFriendRequest_InvalidJSON_Unit(t *testing.T) {
	// Given: An invalid JSON request body
	mockFriendRequestRepo := new(mocks.MockFriendRequestRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewFriendRequestHandlerWithRepo(mockFriendRequestRepo, mockUserRepo)

	app := fiber.New()
	app.Post("/friend-requests", handler.CreateFriendRequest)

	req := httptest.NewRequest("POST", "/friend-requests", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create friend request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestFriendRequestHandler_CreateFriendRequest_MissingSenderID_Unit(t *testing.T) {
	// Given: A request without sender ID
	mockFriendRequestRepo := new(mocks.MockFriendRequestRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewFriendRequestHandlerWithRepo(mockFriendRequestRepo, mockUserRepo)

	app := fiber.New()
	app.Post("/friend-requests", handler.CreateFriendRequest)

	reqBody := dtos.CreateFriendRequestRequest{
		ReceiverID: 2,
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

func TestFriendRequestHandler_CreateFriendRequest_MissingReceiverID_Unit(t *testing.T) {
	// Given: A request without receiver ID
	mockFriendRequestRepo := new(mocks.MockFriendRequestRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewFriendRequestHandlerWithRepo(mockFriendRequestRepo, mockUserRepo)

	app := fiber.New()
	app.Post("/friend-requests", handler.CreateFriendRequest)

	reqBody := dtos.CreateFriendRequestRequest{
		SenderID: 1,
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

func TestFriendRequestHandler_CreateFriendRequest_SelfRequest_Unit(t *testing.T) {
	// Given: A request to send friend request to self
	mockFriendRequestRepo := new(mocks.MockFriendRequestRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewFriendRequestHandlerWithRepo(mockFriendRequestRepo, mockUserRepo)

	app := fiber.New()
	app.Post("/friend-requests", handler.CreateFriendRequest)

	reqBody := dtos.CreateFriendRequestRequest{
		SenderID:   1,
		ReceiverID: 1,
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

func TestFriendRequestHandler_CreateFriendRequest_SenderNotFound_Unit(t *testing.T) {
	// Given: The sender does not exist
	mockFriendRequestRepo := new(mocks.MockFriendRequestRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewFriendRequestHandlerWithRepo(mockFriendRequestRepo, mockUserRepo)

	app := fiber.New()
	app.Post("/friend-requests", handler.CreateFriendRequest)

	mockUserRepo.On("FindByID", mock.Anything, uint(999)).Return(nil, gorm.ErrRecordNotFound)

	reqBody := dtos.CreateFriendRequestRequest{
		SenderID:   999,
		ReceiverID: 2,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/friend-requests", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create friend request
	resp, err := app.Test(req)

	// Then: The request should fail with not found
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

func TestFriendRequestHandler_CreateFriendRequest_ReceiverNotFound_Unit(t *testing.T) {
	// Given: The receiver does not exist
	mockFriendRequestRepo := new(mocks.MockFriendRequestRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewFriendRequestHandlerWithRepo(mockFriendRequestRepo, mockUserRepo)

	app := fiber.New()
	app.Post("/friend-requests", handler.CreateFriendRequest)

	sender := &models.User{Model: gorm.Model{ID: 1}, FirstName: "John"}
	mockUserRepo.On("FindByID", mock.Anything, uint(1)).Return(sender, nil)
	mockUserRepo.On("FindByID", mock.Anything, uint(999)).Return(nil, gorm.ErrRecordNotFound)

	reqBody := dtos.CreateFriendRequestRequest{
		SenderID:   1,
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
	mockUserRepo.AssertExpectations(t)
}

func TestFriendRequestHandler_CreateFriendRequest_AlreadyPending_Unit(t *testing.T) {
	// Given: A pending friend request already exists
	mockFriendRequestRepo := new(mocks.MockFriendRequestRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewFriendRequestHandlerWithRepo(mockFriendRequestRepo, mockUserRepo)

	app := fiber.New()
	app.Post("/friend-requests", handler.CreateFriendRequest)

	sender := &models.User{Model: gorm.Model{ID: 1}, FirstName: "John"}
	receiver := &models.User{Model: gorm.Model{ID: 2}, FirstName: "Jane"}
	mockUserRepo.On("FindByID", mock.Anything, uint(1)).Return(sender, nil)
	mockUserRepo.On("FindByID", mock.Anything, uint(2)).Return(receiver, nil)
	mockFriendRequestRepo.On("FindByUsers", mock.Anything, uint(1), uint(2)).Return(&models.FriendRequest{
		Status: models.StatusPending,
	}, nil)

	reqBody := dtos.CreateFriendRequestRequest{
		SenderID:   1,
		ReceiverID: 2,
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

func TestFriendRequestHandler_CreateFriendRequest_AlreadyFriends_Unit(t *testing.T) {
	// Given: Users are already friends
	mockFriendRequestRepo := new(mocks.MockFriendRequestRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewFriendRequestHandlerWithRepo(mockFriendRequestRepo, mockUserRepo)

	app := fiber.New()
	app.Post("/friend-requests", handler.CreateFriendRequest)

	sender := &models.User{Model: gorm.Model{ID: 1}, FirstName: "John"}
	receiver := &models.User{Model: gorm.Model{ID: 2}, FirstName: "Jane"}
	mockUserRepo.On("FindByID", mock.Anything, uint(1)).Return(sender, nil)
	mockUserRepo.On("FindByID", mock.Anything, uint(2)).Return(receiver, nil)
	mockFriendRequestRepo.On("FindByUsers", mock.Anything, uint(1), uint(2)).Return(&models.FriendRequest{
		Status: models.StatusAccepted,
	}, nil)

	reqBody := dtos.CreateFriendRequestRequest{
		SenderID:   1,
		ReceiverID: 2,
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

func TestFriendRequestHandler_CreateFriendRequest_DatabaseError_Unit(t *testing.T) {
	// Given: A database error occurs when creating friend request
	mockFriendRequestRepo := new(mocks.MockFriendRequestRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewFriendRequestHandlerWithRepo(mockFriendRequestRepo, mockUserRepo)

	app := fiber.New()
	app.Post("/friend-requests", handler.CreateFriendRequest)

	sender := &models.User{Model: gorm.Model{ID: 1}, FirstName: "John"}
	receiver := &models.User{Model: gorm.Model{ID: 2}, FirstName: "Jane"}
	mockUserRepo.On("FindByID", mock.Anything, uint(1)).Return(sender, nil)
	mockUserRepo.On("FindByID", mock.Anything, uint(2)).Return(receiver, nil)
	mockFriendRequestRepo.On("FindByUsers", mock.Anything, uint(1), uint(2)).Return(nil, gorm.ErrRecordNotFound)
	mockFriendRequestRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.FriendRequest")).Return(errors.New("database error"))

	reqBody := dtos.CreateFriendRequestRequest{
		SenderID:   1,
		ReceiverID: 2,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/friend-requests", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create friend request
	resp, err := app.Test(req)

	// Then: The request should fail with internal server error
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockFriendRequestRepo.AssertExpectations(t)
}

func TestFriendRequestHandler_AcceptFriendRequest_Success_Unit(t *testing.T) {
	// Given: A pending friend request exists
	mockFriendRequestRepo := new(mocks.MockFriendRequestRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewFriendRequestHandlerWithRepo(mockFriendRequestRepo, mockUserRepo)

	app := fiber.New()
	app.Put("/friend-requests/:id/accept", handler.AcceptFriendRequest)

	sender := models.User{Model: gorm.Model{ID: 1}, FirstName: "John", LastName: "Doe"}
	receiver := models.User{Model: gorm.Model{ID: 2}, FirstName: "Jane", LastName: "Smith"}
	friendRequest := &models.FriendRequest{
		Model:      gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		SenderID:   1,
		ReceiverID: 2,
		Status:     models.StatusPending,
		Sender:     sender,
		Receiver:   receiver,
	}
	mockFriendRequestRepo.On("FindByID", mock.Anything, uint(1)).Return(friendRequest, nil).Times(2)
	mockFriendRequestRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.FriendRequest")).Return(nil)

	req := httptest.NewRequest("PUT", "/friend-requests/1/accept", nil)

	// When: Making the accept friend request
	resp, err := app.Test(req)

	// Then: The request should succeed
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockFriendRequestRepo.AssertExpectations(t)
}

func TestFriendRequestHandler_AcceptFriendRequest_NotFound_Unit(t *testing.T) {
	// Given: No friend request exists with the specified ID
	mockFriendRequestRepo := new(mocks.MockFriendRequestRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewFriendRequestHandlerWithRepo(mockFriendRequestRepo, mockUserRepo)

	app := fiber.New()
	app.Put("/friend-requests/:id/accept", handler.AcceptFriendRequest)

	mockFriendRequestRepo.On("FindByID", mock.Anything, uint(999)).Return(nil, gorm.ErrRecordNotFound)

	req := httptest.NewRequest("PUT", "/friend-requests/999/accept", nil)

	// When: Making the accept friend request
	resp, err := app.Test(req)

	// Then: The request should fail with not found
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockFriendRequestRepo.AssertExpectations(t)
}

func TestFriendRequestHandler_AcceptFriendRequest_InvalidID_Unit(t *testing.T) {
	// Given: An invalid friend request ID
	mockFriendRequestRepo := new(mocks.MockFriendRequestRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewFriendRequestHandlerWithRepo(mockFriendRequestRepo, mockUserRepo)

	app := fiber.New()
	app.Put("/friend-requests/:id/accept", handler.AcceptFriendRequest)

	req := httptest.NewRequest("PUT", "/friend-requests/invalid/accept", nil)

	// When: Making the accept friend request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestFriendRequestHandler_AcceptFriendRequest_NotPending_Unit(t *testing.T) {
	// Given: A friend request that is not pending
	mockFriendRequestRepo := new(mocks.MockFriendRequestRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewFriendRequestHandlerWithRepo(mockFriendRequestRepo, mockUserRepo)

	app := fiber.New()
	app.Put("/friend-requests/:id/accept", handler.AcceptFriendRequest)

	friendRequest := &models.FriendRequest{
		Model:  gorm.Model{ID: 1},
		Status: models.StatusAccepted,
	}
	mockFriendRequestRepo.On("FindByID", mock.Anything, uint(1)).Return(friendRequest, nil)

	req := httptest.NewRequest("PUT", "/friend-requests/1/accept", nil)

	// When: Making the accept friend request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestFriendRequestHandler_AcceptFriendRequest_DatabaseError_Unit(t *testing.T) {
	// Given: A database error occurs during update
	mockFriendRequestRepo := new(mocks.MockFriendRequestRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewFriendRequestHandlerWithRepo(mockFriendRequestRepo, mockUserRepo)

	app := fiber.New()
	app.Put("/friend-requests/:id/accept", handler.AcceptFriendRequest)

	friendRequest := &models.FriendRequest{
		Model:  gorm.Model{ID: 1},
		Status: models.StatusPending,
	}
	mockFriendRequestRepo.On("FindByID", mock.Anything, uint(1)).Return(friendRequest, nil)
	mockFriendRequestRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.FriendRequest")).Return(errors.New("database error"))

	req := httptest.NewRequest("PUT", "/friend-requests/1/accept", nil)

	// When: Making the accept friend request
	resp, err := app.Test(req)

	// Then: The request should fail with internal server error
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockFriendRequestRepo.AssertExpectations(t)
}

func TestFriendRequestHandler_DeclineFriendRequest_Success_Unit(t *testing.T) {
	// Given: A pending friend request exists
	mockFriendRequestRepo := new(mocks.MockFriendRequestRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewFriendRequestHandlerWithRepo(mockFriendRequestRepo, mockUserRepo)

	app := fiber.New()
	app.Put("/friend-requests/:id/decline", handler.DeclineFriendRequest)

	sender := models.User{Model: gorm.Model{ID: 1}, FirstName: "John", LastName: "Doe"}
	receiver := models.User{Model: gorm.Model{ID: 2}, FirstName: "Jane", LastName: "Smith"}
	friendRequest := &models.FriendRequest{
		Model:      gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		SenderID:   1,
		ReceiverID: 2,
		Status:     models.StatusPending,
		Sender:     sender,
		Receiver:   receiver,
	}
	mockFriendRequestRepo.On("FindByID", mock.Anything, uint(1)).Return(friendRequest, nil).Times(2)
	mockFriendRequestRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.FriendRequest")).Return(nil)

	req := httptest.NewRequest("PUT", "/friend-requests/1/decline", nil)

	// When: Making the decline friend request
	resp, err := app.Test(req)

	// Then: The request should succeed
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockFriendRequestRepo.AssertExpectations(t)
}

func TestFriendRequestHandler_DeclineFriendRequest_NotFound_Unit(t *testing.T) {
	// Given: No friend request exists with the specified ID
	mockFriendRequestRepo := new(mocks.MockFriendRequestRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewFriendRequestHandlerWithRepo(mockFriendRequestRepo, mockUserRepo)

	app := fiber.New()
	app.Put("/friend-requests/:id/decline", handler.DeclineFriendRequest)

	mockFriendRequestRepo.On("FindByID", mock.Anything, uint(999)).Return(nil, gorm.ErrRecordNotFound)

	req := httptest.NewRequest("PUT", "/friend-requests/999/decline", nil)

	// When: Making the decline friend request
	resp, err := app.Test(req)

	// Then: The request should fail with not found
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockFriendRequestRepo.AssertExpectations(t)
}

func TestFriendRequestHandler_DeclineFriendRequest_InvalidID_Unit(t *testing.T) {
	// Given: An invalid friend request ID
	mockFriendRequestRepo := new(mocks.MockFriendRequestRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewFriendRequestHandlerWithRepo(mockFriendRequestRepo, mockUserRepo)

	app := fiber.New()
	app.Put("/friend-requests/:id/decline", handler.DeclineFriendRequest)

	req := httptest.NewRequest("PUT", "/friend-requests/invalid/decline", nil)

	// When: Making the decline friend request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestFriendRequestHandler_DeclineFriendRequest_NotPending_Unit(t *testing.T) {
	// Given: A friend request that is not pending
	mockFriendRequestRepo := new(mocks.MockFriendRequestRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewFriendRequestHandlerWithRepo(mockFriendRequestRepo, mockUserRepo)

	app := fiber.New()
	app.Put("/friend-requests/:id/decline", handler.DeclineFriendRequest)

	friendRequest := &models.FriendRequest{
		Model:  gorm.Model{ID: 1},
		Status: models.StatusDeclined,
	}
	mockFriendRequestRepo.On("FindByID", mock.Anything, uint(1)).Return(friendRequest, nil)

	req := httptest.NewRequest("PUT", "/friend-requests/1/decline", nil)

	// When: Making the decline friend request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestFriendRequestHandler_DeleteFriendRequest_Success_Unit(t *testing.T) {
	// Given: A friend request exists
	mockFriendRequestRepo := new(mocks.MockFriendRequestRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewFriendRequestHandlerWithRepo(mockFriendRequestRepo, mockUserRepo)

	app := fiber.New()
	app.Delete("/friend-requests/:id", handler.DeleteFriendRequest)

	friendRequest := &models.FriendRequest{Model: gorm.Model{ID: 1}}
	mockFriendRequestRepo.On("FindByID", mock.Anything, uint(1)).Return(friendRequest, nil)
	mockFriendRequestRepo.On("Delete", mock.Anything, uint(1)).Return(nil)

	req := httptest.NewRequest("DELETE", "/friend-requests/1", nil)

	// When: Making the delete friend request
	resp, err := app.Test(req)

	// Then: The request should succeed with no content
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	mockFriendRequestRepo.AssertExpectations(t)
}

func TestFriendRequestHandler_DeleteFriendRequest_NotFound_Unit(t *testing.T) {
	// Given: No friend request exists with the specified ID
	mockFriendRequestRepo := new(mocks.MockFriendRequestRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewFriendRequestHandlerWithRepo(mockFriendRequestRepo, mockUserRepo)

	app := fiber.New()
	app.Delete("/friend-requests/:id", handler.DeleteFriendRequest)

	mockFriendRequestRepo.On("FindByID", mock.Anything, uint(999)).Return(nil, gorm.ErrRecordNotFound)

	req := httptest.NewRequest("DELETE", "/friend-requests/999", nil)

	// When: Making the delete friend request
	resp, err := app.Test(req)

	// Then: The request should fail with not found
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockFriendRequestRepo.AssertExpectations(t)
}

func TestFriendRequestHandler_DeleteFriendRequest_InvalidID_Unit(t *testing.T) {
	// Given: An invalid friend request ID
	mockFriendRequestRepo := new(mocks.MockFriendRequestRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewFriendRequestHandlerWithRepo(mockFriendRequestRepo, mockUserRepo)

	app := fiber.New()
	app.Delete("/friend-requests/:id", handler.DeleteFriendRequest)

	req := httptest.NewRequest("DELETE", "/friend-requests/invalid", nil)

	// When: Making the delete friend request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestFriendRequestHandler_DeleteFriendRequest_DatabaseError_Unit(t *testing.T) {
	// Given: A database error occurs during delete
	mockFriendRequestRepo := new(mocks.MockFriendRequestRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewFriendRequestHandlerWithRepo(mockFriendRequestRepo, mockUserRepo)

	app := fiber.New()
	app.Delete("/friend-requests/:id", handler.DeleteFriendRequest)

	friendRequest := &models.FriendRequest{Model: gorm.Model{ID: 1}}
	mockFriendRequestRepo.On("FindByID", mock.Anything, uint(1)).Return(friendRequest, nil)
	mockFriendRequestRepo.On("Delete", mock.Anything, uint(1)).Return(errors.New("database error"))

	req := httptest.NewRequest("DELETE", "/friend-requests/1", nil)

	// When: Making the delete friend request
	resp, err := app.Test(req)

	// Then: The request should fail with internal server error
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockFriendRequestRepo.AssertExpectations(t)
}

func TestFriendRequestHandler_GetSentFriendRequests_Success_Unit(t *testing.T) {
	// Given: A user exists with sent friend requests
	mockFriendRequestRepo := new(mocks.MockFriendRequestRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewFriendRequestHandlerWithRepo(mockFriendRequestRepo, mockUserRepo)

	app := fiber.New()
	app.Get("/users/:id/friend-requests/sent", handler.GetSentFriendRequests)

	user := &models.User{Model: gorm.Model{ID: 1}, FirstName: "John"}
	mockUserRepo.On("FindByID", mock.Anything, uint(1)).Return(user, nil)

	friendRequests := []models.FriendRequest{
		{Model: gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}, SenderID: 1, ReceiverID: 2, Status: models.StatusPending, Sender: *user, Receiver: models.User{FirstName: "Jane"}},
		{Model: gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()}, SenderID: 1, ReceiverID: 3, Status: models.StatusAccepted, Sender: *user, Receiver: models.User{FirstName: "Bob"}},
	}
	mockFriendRequestRepo.On("FindBySenderID", mock.Anything, uint(1)).Return(friendRequests, nil)

	req := httptest.NewRequest("GET", "/users/1/friend-requests/sent", nil)

	// When: Making the get sent friend requests request
	resp, err := app.Test(req)

	// Then: The request should succeed
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
	mockFriendRequestRepo.AssertExpectations(t)
}

func TestFriendRequestHandler_GetSentFriendRequests_UserNotFound_Unit(t *testing.T) {
	// Given: The user does not exist
	mockFriendRequestRepo := new(mocks.MockFriendRequestRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewFriendRequestHandlerWithRepo(mockFriendRequestRepo, mockUserRepo)

	app := fiber.New()
	app.Get("/users/:id/friend-requests/sent", handler.GetSentFriendRequests)

	mockUserRepo.On("FindByID", mock.Anything, uint(999)).Return(nil, gorm.ErrRecordNotFound)

	req := httptest.NewRequest("GET", "/users/999/friend-requests/sent", nil)

	// When: Making the get sent friend requests request
	resp, err := app.Test(req)

	// Then: The request should fail with not found
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

func TestFriendRequestHandler_GetSentFriendRequests_InvalidID_Unit(t *testing.T) {
	// Given: An invalid user ID
	mockFriendRequestRepo := new(mocks.MockFriendRequestRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewFriendRequestHandlerWithRepo(mockFriendRequestRepo, mockUserRepo)

	app := fiber.New()
	app.Get("/users/:id/friend-requests/sent", handler.GetSentFriendRequests)

	req := httptest.NewRequest("GET", "/users/invalid/friend-requests/sent", nil)

	// When: Making the get sent friend requests request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestFriendRequestHandler_GetReceivedFriendRequests_Success_Unit(t *testing.T) {
	// Given: A user exists with received friend requests
	mockFriendRequestRepo := new(mocks.MockFriendRequestRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewFriendRequestHandlerWithRepo(mockFriendRequestRepo, mockUserRepo)

	app := fiber.New()
	app.Get("/users/:id/friend-requests/received", handler.GetReceivedFriendRequests)

	user := &models.User{Model: gorm.Model{ID: 1}, FirstName: "John"}
	mockUserRepo.On("FindByID", mock.Anything, uint(1)).Return(user, nil)

	friendRequests := []models.FriendRequest{
		{Model: gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}, SenderID: 2, ReceiverID: 1, Status: models.StatusPending, Sender: models.User{FirstName: "Jane"}, Receiver: *user},
	}
	mockFriendRequestRepo.On("FindPendingByReceiverID", mock.Anything, uint(1)).Return(friendRequests, nil)

	req := httptest.NewRequest("GET", "/users/1/friend-requests/received", nil)

	// When: Making the get received friend requests request
	resp, err := app.Test(req)

	// Then: The request should succeed
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
	mockFriendRequestRepo.AssertExpectations(t)
}

func TestFriendRequestHandler_GetReceivedFriendRequests_UserNotFound_Unit(t *testing.T) {
	// Given: The user does not exist
	mockFriendRequestRepo := new(mocks.MockFriendRequestRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewFriendRequestHandlerWithRepo(mockFriendRequestRepo, mockUserRepo)

	app := fiber.New()
	app.Get("/users/:id/friend-requests/received", handler.GetReceivedFriendRequests)

	mockUserRepo.On("FindByID", mock.Anything, uint(999)).Return(nil, gorm.ErrRecordNotFound)

	req := httptest.NewRequest("GET", "/users/999/friend-requests/received", nil)

	// When: Making the get received friend requests request
	resp, err := app.Test(req)

	// Then: The request should fail with not found
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

func TestFriendRequestHandler_GetReceivedFriendRequests_InvalidID_Unit(t *testing.T) {
	// Given: An invalid user ID
	mockFriendRequestRepo := new(mocks.MockFriendRequestRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewFriendRequestHandlerWithRepo(mockFriendRequestRepo, mockUserRepo)

	app := fiber.New()
	app.Get("/users/:id/friend-requests/received", handler.GetReceivedFriendRequests)

	req := httptest.NewRequest("GET", "/users/invalid/friend-requests/received", nil)

	// When: Making the get received friend requests request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestFriendRequestHandler_GetFriends_Success_Unit(t *testing.T) {
	// Given: A user exists with friends
	mockFriendRequestRepo := new(mocks.MockFriendRequestRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewFriendRequestHandlerWithRepo(mockFriendRequestRepo, mockUserRepo)

	app := fiber.New()
	app.Get("/users/:id/friends", handler.GetFriends)

	user := &models.User{Model: gorm.Model{ID: 1}, FirstName: "John"}
	mockUserRepo.On("FindByID", mock.Anything, uint(1)).Return(user, nil)

	friends := []models.User{
		{Model: gorm.Model{ID: 2}, FirstName: "Jane", LastName: "Doe", Email: "jane@test.com"},
		{Model: gorm.Model{ID: 3}, FirstName: "Bob", LastName: "Smith", Email: "bob@test.com"},
	}
	mockFriendRequestRepo.On("FindFriendsByUserID", mock.Anything, uint(1)).Return(friends, nil)

	req := httptest.NewRequest("GET", "/users/1/friends", nil)

	// When: Making the get friends request
	resp, err := app.Test(req)

	// Then: The request should succeed
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
	mockFriendRequestRepo.AssertExpectations(t)
}

func TestFriendRequestHandler_GetFriends_UserNotFound_Unit(t *testing.T) {
	// Given: The user does not exist
	mockFriendRequestRepo := new(mocks.MockFriendRequestRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewFriendRequestHandlerWithRepo(mockFriendRequestRepo, mockUserRepo)

	app := fiber.New()
	app.Get("/users/:id/friends", handler.GetFriends)

	mockUserRepo.On("FindByID", mock.Anything, uint(999)).Return(nil, gorm.ErrRecordNotFound)

	req := httptest.NewRequest("GET", "/users/999/friends", nil)

	// When: Making the get friends request
	resp, err := app.Test(req)

	// Then: The request should fail with not found
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

func TestFriendRequestHandler_GetFriends_InvalidID_Unit(t *testing.T) {
	// Given: An invalid user ID
	mockFriendRequestRepo := new(mocks.MockFriendRequestRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewFriendRequestHandlerWithRepo(mockFriendRequestRepo, mockUserRepo)

	app := fiber.New()
	app.Get("/users/:id/friends", handler.GetFriends)

	req := httptest.NewRequest("GET", "/users/invalid/friends", nil)

	// When: Making the get friends request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestFriendRequestHandler_GetFriends_DatabaseError_Unit(t *testing.T) {
	// Given: A database error occurs
	mockFriendRequestRepo := new(mocks.MockFriendRequestRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewFriendRequestHandlerWithRepo(mockFriendRequestRepo, mockUserRepo)

	app := fiber.New()
	app.Get("/users/:id/friends", handler.GetFriends)

	user := &models.User{Model: gorm.Model{ID: 1}, FirstName: "John"}
	mockUserRepo.On("FindByID", mock.Anything, uint(1)).Return(user, nil)
	mockFriendRequestRepo.On("FindFriendsByUserID", mock.Anything, uint(1)).Return(nil, errors.New("database error"))

	req := httptest.NewRequest("GET", "/users/1/friends", nil)

	// When: Making the get friends request
	resp, err := app.Test(req)

	// Then: The request should fail with internal server error
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockFriendRequestRepo.AssertExpectations(t)
}

func TestNewFriendRequestHandlerWithRepo_Unit(t *testing.T) {
	// Given: Mock repositories
	mockFriendRequestRepo := new(mocks.MockFriendRequestRepository)
	mockUserRepo := new(mocks.MockUserRepository)

	// When: Creating a new friend request handler with repos
	handler := NewFriendRequestHandlerWithRepo(mockFriendRequestRepo, mockUserRepo)

	// Then: The handler should be created
	assert.NotNil(t, handler)
}
