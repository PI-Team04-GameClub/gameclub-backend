package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/mocks"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestUserHandler_GetAllUsers_Success_Unit(t *testing.T) {
	// Given: Users exist in the repository
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewUserHandlerWithRepo(mockUserRepo)

	app := fiber.New()
	app.Get("/users", handler.GetAllUsers)

	users := []models.User{
		{Model: gorm.Model{ID: 1}, FirstName: "John", LastName: "Doe", Email: "john@example.com"},
		{Model: gorm.Model{ID: 2}, FirstName: "Jane", LastName: "Doe", Email: "jane@example.com"},
	}
	mockUserRepo.On("FindAll", mock.Anything).Return(users, nil)

	req := httptest.NewRequest("GET", "/users", nil)

	// When: Making the get all users request
	resp, err := app.Test(req)

	// Then: The request should succeed with users list
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

func TestUserHandler_GetAllUsers_Empty_Unit(t *testing.T) {
	// Given: No users exist in the repository
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewUserHandlerWithRepo(mockUserRepo)

	app := fiber.New()
	app.Get("/users", handler.GetAllUsers)

	mockUserRepo.On("FindAll", mock.Anything).Return([]models.User{}, nil)

	req := httptest.NewRequest("GET", "/users", nil)

	// When: Making the get all users request
	resp, err := app.Test(req)

	// Then: The request should succeed with empty list
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

func TestUserHandler_GetAllUsers_DatabaseError_Unit(t *testing.T) {
	// Given: A database error occurs
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewUserHandlerWithRepo(mockUserRepo)

	app := fiber.New()
	app.Get("/users", handler.GetAllUsers)

	mockUserRepo.On("FindAll", mock.Anything).Return(nil, errors.New("database error"))

	req := httptest.NewRequest("GET", "/users", nil)

	// When: Making the get all users request
	resp, err := app.Test(req)

	// Then: The request should fail with internal server error
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

func TestUserHandler_GetUserByID_Success_Unit(t *testing.T) {
	// Given: A user exists with the specified ID
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewUserHandlerWithRepo(mockUserRepo)

	app := fiber.New()
	app.Get("/users/:id", handler.GetUserByID)

	user := &models.User{Model: gorm.Model{ID: 1}, FirstName: "John", LastName: "Doe", Email: "john@example.com"}
	mockUserRepo.On("FindByID", mock.Anything, uint(1)).Return(user, nil)

	req := httptest.NewRequest("GET", "/users/1", nil)

	// When: Making the get user by ID request
	resp, err := app.Test(req)

	// Then: The request should succeed with user data
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

func TestUserHandler_GetUserByID_NotFound_Unit(t *testing.T) {
	// Given: No user exists with the specified ID
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewUserHandlerWithRepo(mockUserRepo)

	app := fiber.New()
	app.Get("/users/:id", handler.GetUserByID)

	mockUserRepo.On("FindByID", mock.Anything, uint(999)).Return(nil, gorm.ErrRecordNotFound)

	req := httptest.NewRequest("GET", "/users/999", nil)

	// When: Making the get user by ID request
	resp, err := app.Test(req)

	// Then: The request should fail with not found status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

func TestUserHandler_GetUserByID_InvalidID_Unit(t *testing.T) {
	// Given: An invalid user ID is provided
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewUserHandlerWithRepo(mockUserRepo)

	app := fiber.New()
	app.Get("/users/:id", handler.GetUserByID)

	req := httptest.NewRequest("GET", "/users/invalid", nil)

	// When: Making the get user by ID request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestUserHandler_UpdateUser_Success_Unit(t *testing.T) {
	// Given: A user exists and valid update request
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewUserHandlerWithRepo(mockUserRepo)

	app := fiber.New()
	app.Put("/users/:id", handler.UpdateUser)

	existingUser := &models.User{Model: gorm.Model{ID: 1}, FirstName: "John", LastName: "Doe", Email: "john@example.com"}
	mockUserRepo.On("FindByID", mock.Anything, uint(1)).Return(existingUser, nil)
	mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)

	reqBody := dtos.UpdateUserRequest{
		FirstName: "Johnny",
		LastName:  "Updated",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/users/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update user request
	resp, err := app.Test(req)

	// Then: The request should succeed
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

func TestUserHandler_UpdateUser_NotFound_Unit(t *testing.T) {
	// Given: No user exists with the specified ID
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewUserHandlerWithRepo(mockUserRepo)

	app := fiber.New()
	app.Put("/users/:id", handler.UpdateUser)

	mockUserRepo.On("FindByID", mock.Anything, uint(999)).Return(nil, gorm.ErrRecordNotFound)

	reqBody := dtos.UpdateUserRequest{
		FirstName: "Johnny",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/users/999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update user request
	resp, err := app.Test(req)

	// Then: The request should fail with not found status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

func TestUserHandler_UpdateUser_InvalidID_Unit(t *testing.T) {
	// Given: An invalid user ID is provided
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewUserHandlerWithRepo(mockUserRepo)

	app := fiber.New()
	app.Put("/users/:id", handler.UpdateUser)

	reqBody := dtos.UpdateUserRequest{
		FirstName: "Johnny",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/users/invalid", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update user request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestUserHandler_UpdateUser_InvalidJSON_Unit(t *testing.T) {
	// Given: A user exists but request body is invalid JSON
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewUserHandlerWithRepo(mockUserRepo)

	app := fiber.New()
	app.Put("/users/:id", handler.UpdateUser)

	existingUser := &models.User{Model: gorm.Model{ID: 1}, FirstName: "John", Email: "john@example.com"}
	mockUserRepo.On("FindByID", mock.Anything, uint(1)).Return(existingUser, nil)

	req := httptest.NewRequest("PUT", "/users/1", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update user request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

func TestUserHandler_UpdateUser_EmailAlreadyExists_Unit(t *testing.T) {
	// Given: A user exists but new email is already taken
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewUserHandlerWithRepo(mockUserRepo)

	app := fiber.New()
	app.Put("/users/:id", handler.UpdateUser)

	existingUser := &models.User{Model: gorm.Model{ID: 1}, FirstName: "John", Email: "john@example.com"}
	otherUser := &models.User{Model: gorm.Model{ID: 2}, Email: "taken@example.com"}
	mockUserRepo.On("FindByID", mock.Anything, uint(1)).Return(existingUser, nil)
	mockUserRepo.On("FindByEmail", mock.Anything, "taken@example.com").Return(otherUser, nil)

	reqBody := dtos.UpdateUserRequest{
		Email: "taken@example.com",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/users/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update user request
	resp, err := app.Test(req)

	// Then: The request should fail with conflict status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusConflict, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

func TestUserHandler_UpdateUser_PasswordTooShort_Unit(t *testing.T) {
	// Given: A user exists but new password is too short
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewUserHandlerWithRepo(mockUserRepo)

	app := fiber.New()
	app.Put("/users/:id", handler.UpdateUser)

	existingUser := &models.User{Model: gorm.Model{ID: 1}, FirstName: "John", Email: "john@example.com"}
	mockUserRepo.On("FindByID", mock.Anything, uint(1)).Return(existingUser, nil)

	reqBody := dtos.UpdateUserRequest{
		Password: "12345",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/users/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update user request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

func TestUserHandler_UpdateUser_DatabaseError_Unit(t *testing.T) {
	// Given: A user exists but database error occurs during update
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewUserHandlerWithRepo(mockUserRepo)

	app := fiber.New()
	app.Put("/users/:id", handler.UpdateUser)

	existingUser := &models.User{Model: gorm.Model{ID: 1}, FirstName: "John", Email: "john@example.com"}
	mockUserRepo.On("FindByID", mock.Anything, uint(1)).Return(existingUser, nil)
	mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.User")).Return(errors.New("database error"))

	reqBody := dtos.UpdateUserRequest{
		FirstName: "Johnny",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/users/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update user request
	resp, err := app.Test(req)

	// Then: The request should fail with internal server error
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}
