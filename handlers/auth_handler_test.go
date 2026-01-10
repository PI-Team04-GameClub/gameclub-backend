package handlers

import (
	"bytes"
	"encoding/json"
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

func TestAuthHandler_Register_Success_Unit(t *testing.T) {
	// Given: A valid registration request and mock repository
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewAuthHandlerWithRepo(mockUserRepo)

	app := fiber.New()
	app.Post("/auth/register", handler.Register)

	mockUserRepo.On("FindByEmail", mock.Anything, "john@example.com").Return(nil, gorm.ErrRecordNotFound)
	mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)

	reqBody := dtos.RegisterRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Password:  "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the registration request
	resp, err := app.Test(req)

	// Then: The registration should succeed
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

func TestAuthHandler_Register_MissingEmail_Unit(t *testing.T) {
	// Given: A registration request without email
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewAuthHandlerWithRepo(mockUserRepo)

	app := fiber.New()
	app.Post("/auth/register", handler.Register)

	reqBody := dtos.RegisterRequest{
		FirstName: "John",
		LastName:  "Doe",
		Password:  "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the registration request
	resp, err := app.Test(req)

	// Then: The registration should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestAuthHandler_Register_MissingPassword_Unit(t *testing.T) {
	// Given: A registration request without password
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewAuthHandlerWithRepo(mockUserRepo)

	app := fiber.New()
	app.Post("/auth/register", handler.Register)

	reqBody := dtos.RegisterRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the registration request
	resp, err := app.Test(req)

	// Then: The registration should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestAuthHandler_Register_MissingFirstName_Unit(t *testing.T) {
	// Given: A registration request without first name
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewAuthHandlerWithRepo(mockUserRepo)

	app := fiber.New()
	app.Post("/auth/register", handler.Register)

	reqBody := dtos.RegisterRequest{
		LastName: "Doe",
		Email:    "john@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the registration request
	resp, err := app.Test(req)

	// Then: The registration should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestAuthHandler_Register_ShortPassword_Unit(t *testing.T) {
	// Given: A registration request with password less than 6 characters
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewAuthHandlerWithRepo(mockUserRepo)

	app := fiber.New()
	app.Post("/auth/register", handler.Register)

	reqBody := dtos.RegisterRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Password:  "12345",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the registration request
	resp, err := app.Test(req)

	// Then: The registration should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestAuthHandler_Register_DuplicateEmail_Unit(t *testing.T) {
	// Given: A user already exists with the same email
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewAuthHandlerWithRepo(mockUserRepo)

	app := fiber.New()
	app.Post("/auth/register", handler.Register)

	existingUser := &models.User{
		Model:     gorm.Model{ID: 1},
		FirstName: "Existing",
		LastName:  "User",
		Email:     "existing@example.com",
	}
	mockUserRepo.On("FindByEmail", mock.Anything, "existing@example.com").Return(existingUser, nil)

	reqBody := dtos.RegisterRequest{
		FirstName: "New",
		LastName:  "User",
		Email:     "existing@example.com",
		Password:  "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the registration request with duplicate email
	resp, err := app.Test(req)

	// Then: The registration should fail with conflict status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusConflict, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

func TestAuthHandler_Register_InvalidJSON_Unit(t *testing.T) {
	// Given: An invalid JSON request body
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewAuthHandlerWithRepo(mockUserRepo)

	app := fiber.New()
	app.Post("/auth/register", handler.Register)

	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the registration request
	resp, err := app.Test(req)

	// Then: The registration should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestAuthHandler_Register_DatabaseError_Unit(t *testing.T) {
	// Given: A database error occurs when checking for existing user
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewAuthHandlerWithRepo(mockUserRepo)

	app := fiber.New()
	app.Post("/auth/register", handler.Register)

	mockUserRepo.On("FindByEmail", mock.Anything, "john@example.com").Return(nil, gorm.ErrInvalidDB)

	reqBody := dtos.RegisterRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Password:  "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the registration request
	resp, err := app.Test(req)

	// Then: The registration should fail with internal server error
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

func TestAuthHandler_Register_CreateFails_Unit(t *testing.T) {
	// Given: Creating user in database fails
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewAuthHandlerWithRepo(mockUserRepo)

	app := fiber.New()
	app.Post("/auth/register", handler.Register)

	mockUserRepo.On("FindByEmail", mock.Anything, "john@example.com").Return(nil, gorm.ErrRecordNotFound)
	mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(gorm.ErrInvalidDB)

	reqBody := dtos.RegisterRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Password:  "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the registration request
	resp, err := app.Test(req)

	// Then: The registration should fail with internal server error
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

func TestAuthHandler_Login_Success_Unit(t *testing.T) {
	// Given: An existing user and valid login credentials
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewAuthHandlerWithRepo(mockUserRepo)

	app := fiber.New()
	app.Post("/auth/login", handler.Login)

	user := &models.User{
		Model:     gorm.Model{ID: 1},
		FirstName: "Login",
		LastName:  "User",
		Email:     "login@example.com",
		Password:  hashPassword("password123"),
	}
	mockUserRepo.On("FindByEmail", mock.Anything, "login@example.com").Return(user, nil)

	reqBody := dtos.LoginRequest{
		Email:    "login@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the login request
	resp, err := app.Test(req)

	// Then: The login should succeed
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

func TestAuthHandler_Login_WrongPassword_Unit(t *testing.T) {
	// Given: An existing user and wrong password
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewAuthHandlerWithRepo(mockUserRepo)

	app := fiber.New()
	app.Post("/auth/login", handler.Login)

	user := &models.User{
		Model:     gorm.Model{ID: 1},
		FirstName: "Wrong",
		LastName:  "Password",
		Email:     "wrong@example.com",
		Password:  hashPassword("correctpassword"),
	}
	mockUserRepo.On("FindByEmail", mock.Anything, "wrong@example.com").Return(user, nil)

	reqBody := dtos.LoginRequest{
		Email:    "wrong@example.com",
		Password: "wrongpassword",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the login request with wrong password
	resp, err := app.Test(req)

	// Then: The login should fail with unauthorized status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

func TestAuthHandler_Login_NonexistentUser_Unit(t *testing.T) {
	// Given: A login request for a non-existent user
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewAuthHandlerWithRepo(mockUserRepo)

	app := fiber.New()
	app.Post("/auth/login", handler.Login)

	mockUserRepo.On("FindByEmail", mock.Anything, "nonexistent@example.com").Return(nil, gorm.ErrRecordNotFound)

	reqBody := dtos.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the login request
	resp, err := app.Test(req)

	// Then: The login should fail with unauthorized status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

func TestAuthHandler_Login_MissingEmail_Unit(t *testing.T) {
	// Given: A login request without email
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewAuthHandlerWithRepo(mockUserRepo)

	app := fiber.New()
	app.Post("/auth/login", handler.Login)

	reqBody := dtos.LoginRequest{
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the login request
	resp, err := app.Test(req)

	// Then: The login should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestAuthHandler_Login_MissingPassword_Unit(t *testing.T) {
	// Given: A login request without password
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewAuthHandlerWithRepo(mockUserRepo)

	app := fiber.New()
	app.Post("/auth/login", handler.Login)

	reqBody := dtos.LoginRequest{
		Email: "user@example.com",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the login request
	resp, err := app.Test(req)

	// Then: The login should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestAuthHandler_Login_InvalidJSON_Unit(t *testing.T) {
	// Given: An invalid JSON request body
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewAuthHandlerWithRepo(mockUserRepo)

	app := fiber.New()
	app.Post("/auth/login", handler.Login)

	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the login request
	resp, err := app.Test(req)

	// Then: The login should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestAuthHandler_Login_DatabaseError_Unit(t *testing.T) {
	// Given: A database error occurs when finding user
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewAuthHandlerWithRepo(mockUserRepo)

	app := fiber.New()
	app.Post("/auth/login", handler.Login)

	mockUserRepo.On("FindByEmail", mock.Anything, "user@example.com").Return(nil, gorm.ErrInvalidDB)

	reqBody := dtos.LoginRequest{
		Email:    "user@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the login request
	resp, err := app.Test(req)

	// Then: The login should fail with internal server error
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

func TestAuthHandler_GetCurrentUser_Success_Unit(t *testing.T) {
	// Given: A request with authenticated user in context
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewAuthHandlerWithRepo(mockUserRepo)

	app := fiber.New()
	app.Get("/auth/me", func(c *fiber.Ctx) error {
		user := &models.User{
			Model:     gorm.Model{ID: 1},
			FirstName: "Test",
			LastName:  "User",
			Email:     "test@example.com",
		}
		c.Locals("user", user)
		return handler.GetCurrentUser(c)
	})

	req := httptest.NewRequest("GET", "/auth/me", nil)

	// When: Making the get current user request
	resp, err := app.Test(req)

	// Then: The request should succeed with user data
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Test", response["first_name"])
	assert.Equal(t, "User", response["last_name"])
	assert.Equal(t, "test@example.com", response["email"])
}
