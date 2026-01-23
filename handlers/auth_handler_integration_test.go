package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Team{}, &models.Game{}, &models.Tournament{}, &models.News{}, &models.Comment{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func setupTestApp(db *gorm.DB) *fiber.App {
	app := fiber.New()
	authHandler := NewAuthHandler(db)

	app.Post("/auth/register", authHandler.Register)
	app.Post("/auth/login", authHandler.Login)
	app.Get("/auth/me", func(c *fiber.Ctx) error {
		user := &models.User{
			Model:     gorm.Model{ID: 1},
			FirstName: "Test",
			LastName:  "User",
			Email:     "test@example.com",
		}
		c.Locals("user", user)
		return authHandler.GetCurrentUser(c)
	})

	return app
}

func TestAuthHandler_Register_Success(t *testing.T) {
	// Given: A valid registration request
	db := setupTestDB(t)
	app := setupTestApp(db)

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
}

func TestAuthHandler_Register_MissingEmail(t *testing.T) {
	// Given: A registration request without email
	db := setupTestDB(t)
	app := setupTestApp(db)

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

func TestAuthHandler_Register_MissingPassword(t *testing.T) {
	// Given: A registration request without password
	db := setupTestDB(t)
	app := setupTestApp(db)

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

func TestAuthHandler_Register_MissingFirstName(t *testing.T) {
	// Given: A registration request without first name
	db := setupTestDB(t)
	app := setupTestApp(db)

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

func TestAuthHandler_Register_ShortPassword(t *testing.T) {
	// Given: A registration request with password less than 6 characters
	db := setupTestDB(t)
	app := setupTestApp(db)

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

func TestAuthHandler_Register_DuplicateEmail(t *testing.T) {
	// Given: A user already exists with the same email
	db := setupTestDB(t)
	app := setupTestApp(db)

	existingUser := models.User{
		FirstName: "Existing",
		LastName:  "User",
		Email:     "existing@example.com",
		Password:  hashPassword("password123"),
	}
	db.Create(&existingUser)

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
}

func TestAuthHandler_Register_InvalidJSON(t *testing.T) {
	// Given: An invalid JSON request body
	db := setupTestDB(t)
	app := setupTestApp(db)

	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the registration request
	resp, err := app.Test(req)

	// Then: The registration should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestAuthHandler_Register_ReturnsToken(t *testing.T) {
	// Given: A valid registration request
	db := setupTestDB(t)
	app := setupTestApp(db)

	reqBody := dtos.RegisterRequest{
		FirstName: "Token",
		LastName:  "User",
		Email:     "token@example.com",
		Password:  "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the registration request
	resp, err := app.Test(req)

	// Then: The response should contain a token
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var response dtos.AuthResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.NotEmpty(t, response.Token)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	// Given: An existing user and valid login credentials
	db := setupTestDB(t)
	app := setupTestApp(db)

	user := models.User{
		FirstName: "Login",
		LastName:  "User",
		Email:     "login@example.com",
		Password:  hashPassword("password123"),
	}
	db.Create(&user)

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
}

func TestAuthHandler_Login_WrongPassword(t *testing.T) {
	// Given: An existing user and wrong password
	db := setupTestDB(t)
	app := setupTestApp(db)

	user := models.User{
		FirstName: "Wrong",
		LastName:  "Password",
		Email:     "wrong@example.com",
		Password:  hashPassword("correctpassword"),
	}
	db.Create(&user)

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
}

func TestAuthHandler_Login_NonexistentUser(t *testing.T) {
	// Given: A login request for a non-existent user
	db := setupTestDB(t)
	app := setupTestApp(db)

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
}

func TestAuthHandler_Login_MissingEmail(t *testing.T) {
	// Given: A login request without email
	db := setupTestDB(t)
	app := setupTestApp(db)

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

func TestAuthHandler_Login_MissingPassword(t *testing.T) {
	// Given: A login request without password
	db := setupTestDB(t)
	app := setupTestApp(db)

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

func TestAuthHandler_Login_InvalidJSON(t *testing.T) {
	// Given: An invalid JSON request body
	db := setupTestDB(t)
	app := setupTestApp(db)

	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the login request
	resp, err := app.Test(req)

	// Then: The login should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestAuthHandler_Login_ReturnsToken(t *testing.T) {
	// Given: An existing user and valid credentials
	db := setupTestDB(t)
	app := setupTestApp(db)

	user := models.User{
		FirstName: "Token",
		LastName:  "Login",
		Email:     "tokenlogin@example.com",
		Password:  hashPassword("password123"),
	}
	db.Create(&user)

	reqBody := dtos.LoginRequest{
		Email:    "tokenlogin@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the login request
	resp, err := app.Test(req)

	// Then: The response should contain a token
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response dtos.AuthResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.NotEmpty(t, response.Token)
}

func TestAuthHandler_GetCurrentUser_Success(t *testing.T) {
	// Given: A request with authenticated user in context
	db := setupTestDB(t)
	app := setupTestApp(db)

	req := httptest.NewRequest("GET", "/auth/me", nil)

	// When: Making the get current user request
	resp, err := app.Test(req)

	// Then: The request should succeed
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestAuthHandler_GetCurrentUser_ReturnsUserData(t *testing.T) {
	// Given: A request with authenticated user in context
	db := setupTestDB(t)
	app := setupTestApp(db)

	req := httptest.NewRequest("GET", "/auth/me", nil)

	// When: Making the get current user request
	resp, err := app.Test(req)

	// Then: The response should contain user data
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Test", response["first_name"])
	assert.Equal(t, "User", response["last_name"])
	assert.Equal(t, "test@example.com", response["email"])
}

func TestHashPassword(t *testing.T) {
	// Given: A password string
	password := "testpassword123"

	// When: Hashing the password
	hashed := hashPassword(password)

	// Then: The hash should be non-empty and different from the original
	assert.NotEmpty(t, hashed)
	assert.NotEqual(t, password, hashed)
}

func TestHashPassword_Consistent(t *testing.T) {
	// Given: The same password hashed twice
	password := "samepassword"

	// When: Hashing the password twice
	hash1 := hashPassword(password)
	hash2 := hashPassword(password)

	// Then: The hashes should be identical
	assert.Equal(t, hash1, hash2)
}

func TestVerifyPassword_Correct(t *testing.T) {
	// Given: A hashed password and the correct password
	password := "correctpassword"
	hashed := hashPassword(password)

	// When: Verifying the password
	result := verifyPassword(hashed, password)

	// Then: The verification should succeed
	assert.True(t, result)
}

func TestVerifyPassword_Incorrect(t *testing.T) {
	// Given: A hashed password and an incorrect password
	password := "correctpassword"
	hashed := hashPassword(password)

	// When: Verifying with wrong password
	result := verifyPassword(hashed, "wrongpassword")

	// Then: The verification should fail
	assert.False(t, result)
}
