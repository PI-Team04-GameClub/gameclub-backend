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

func setupUserTestApp(db *gorm.DB) *fiber.App {
	app := fiber.New()
	userHandler := NewUserHandler(db)

	app.Get("/users", userHandler.GetAllUsers)
	app.Get("/users/:id", userHandler.GetUserByID)
	app.Put("/users/:id", userHandler.UpdateUser)

	return app
}

func TestUserHandler_GetAllUsers_Empty(t *testing.T) {
	// Given: An empty database
	db := setupTestDB(t)
	app := setupUserTestApp(db)

	req := httptest.NewRequest("GET", "/users", nil)

	// When: Making the get all users request
	resp, err := app.Test(req)

	// Then: The request should succeed with empty array
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var users []dtos.UserResponse
	json.NewDecoder(resp.Body).Decode(&users)
	assert.Empty(t, users)
}

func TestUserHandler_GetAllUsers_WithUsers(t *testing.T) {
	// Given: A database with users
	db := setupTestDB(t)
	app := setupUserTestApp(db)

	db.Create(&models.User{FirstName: "John", LastName: "Doe", Email: "john@example.com", Password: "hashed"})
	db.Create(&models.User{FirstName: "Jane", LastName: "Doe", Email: "jane@example.com", Password: "hashed"})

	req := httptest.NewRequest("GET", "/users", nil)

	// When: Making the get all users request
	resp, err := app.Test(req)

	// Then: The request should return all users
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var users []dtos.UserResponse
	json.NewDecoder(resp.Body).Decode(&users)
	assert.Len(t, users, 2)
}

func TestUserHandler_GetUserByID_Found(t *testing.T) {
	// Given: A user exists in the database
	db := setupTestDB(t)
	app := setupUserTestApp(db)

	user := models.User{FirstName: "John", LastName: "Doe", Email: "john@example.com", Password: "hashed"}
	db.Create(&user)

	req := httptest.NewRequest("GET", fmt.Sprintf("/users/%d", user.ID), nil)

	// When: Making the get user by ID request
	resp, err := app.Test(req)

	// Then: The request should return the user
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response dtos.UserResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "John", response.FirstName)
	assert.Equal(t, "john@example.com", response.Email)
}

func TestUserHandler_GetUserByID_NotFound(t *testing.T) {
	// Given: An empty database
	db := setupTestDB(t)
	app := setupUserTestApp(db)

	req := httptest.NewRequest("GET", "/users/999", nil)

	// When: Making the get user by ID request for non-existent user
	resp, err := app.Test(req)

	// Then: The request should return not found
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestUserHandler_GetUserByID_InvalidID(t *testing.T) {
	// Given: An invalid user ID
	db := setupTestDB(t)
	app := setupUserTestApp(db)

	req := httptest.NewRequest("GET", "/users/invalid", nil)

	// When: Making the get user by ID request with invalid ID
	resp, err := app.Test(req)

	// Then: The request should return bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestUserHandler_UpdateUser_Success(t *testing.T) {
	// Given: An existing user
	db := setupTestDB(t)
	app := setupUserTestApp(db)

	user := models.User{FirstName: "John", LastName: "Doe", Email: "john@example.com", Password: "hashed"}
	db.Create(&user)

	reqBody := dtos.UpdateUserRequest{
		FirstName: "Johnny",
		LastName:  "Updated",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", fmt.Sprintf("/users/%d", user.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update user request
	resp, err := app.Test(req)

	// Then: The user should be updated
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response dtos.UserResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Johnny", response.FirstName)
	assert.Equal(t, "Updated", response.LastName)
}

func TestUserHandler_UpdateUser_NotFound(t *testing.T) {
	// Given: An empty database
	db := setupTestDB(t)
	app := setupUserTestApp(db)

	reqBody := dtos.UpdateUserRequest{FirstName: "Updated"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/users/999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update request for non-existent user
	resp, err := app.Test(req)

	// Then: The request should return not found
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestUserHandler_UpdateUser_InvalidJSON(t *testing.T) {
	// Given: An existing user and invalid JSON
	db := setupTestDB(t)
	app := setupUserTestApp(db)

	user := models.User{FirstName: "John", Email: "john@example.com", Password: "hashed"}
	db.Create(&user)

	req := httptest.NewRequest("PUT", fmt.Sprintf("/users/%d", user.ID), bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update request with invalid JSON
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestUserHandler_UpdateUser_DuplicateEmail(t *testing.T) {
	// Given: Two existing users
	db := setupTestDB(t)
	app := setupUserTestApp(db)

	user1 := models.User{FirstName: "John", Email: "john@example.com", Password: "hashed"}
	user2 := models.User{FirstName: "Jane", Email: "jane@example.com", Password: "hashed"}
	db.Create(&user1)
	db.Create(&user2)

	reqBody := dtos.UpdateUserRequest{
		Email: "jane@example.com",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", fmt.Sprintf("/users/%d", user1.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update request with duplicate email
	resp, err := app.Test(req)

	// Then: The request should fail with conflict
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusConflict, resp.StatusCode)
}

func TestUserHandler_UpdateUser_PasswordTooShort(t *testing.T) {
	// Given: An existing user
	db := setupTestDB(t)
	app := setupUserTestApp(db)

	user := models.User{FirstName: "John", Email: "john@example.com", Password: "hashed"}
	db.Create(&user)

	reqBody := dtos.UpdateUserRequest{
		Password: "12345",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", fmt.Sprintf("/users/%d", user.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update request with short password
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestUserHandler_UpdateUser_PasswordIsHashed(t *testing.T) {
	// Given: An existing user
	db := setupTestDB(t)
	app := setupUserTestApp(db)

	user := models.User{FirstName: "John", Email: "john@example.com", Password: "oldhash"}
	db.Create(&user)

	reqBody := dtos.UpdateUserRequest{
		Password: "newpassword123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", fmt.Sprintf("/users/%d", user.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update user request
	app.Test(req)

	// Then: The password in the database should be hashed
	var updatedUser models.User
	db.First(&updatedUser, user.ID)
	assert.NotEqual(t, "newpassword123", updatedUser.Password)
	assert.NotEqual(t, "oldhash", updatedUser.Password)
}

func TestNewUserHandler(t *testing.T) {
	// Given: A database connection
	db := setupTestDB(t)

	// When: Creating a new user handler
	handler := NewUserHandler(db)

	// Then: The handler should be created
	assert.NotNil(t, handler)
}
