package middleware

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/PI-Team04-GameClub/gameclub-backend/security"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupMiddlewareTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func setupMiddlewareTestApp(db *gorm.DB) *fiber.App {
	app := fiber.New()

	app.Get("/protected", JWTMiddleware(db), func(c *fiber.Ctx) error {
		user := c.Locals("user").(*models.User)
		return c.JSON(fiber.Map{
			"id":    user.ID,
			"email": user.Email,
		})
	})

	return app
}

func TestJWTMiddleware_MissingAuthorizationHeader(t *testing.T) {
	// Given: A request without Authorization header
	db := setupMiddlewareTestDB(t)
	app := setupMiddlewareTestApp(db)

	req := httptest.NewRequest("GET", "/protected", nil)

	// When: Making the request
	resp, err := app.Test(req)

	// Then: The request should be rejected with unauthorized status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Missing authorization header", response["error"])
}

func TestJWTMiddleware_InvalidToken(t *testing.T) {
	// Given: A request with an invalid token
	db := setupMiddlewareTestDB(t)
	app := setupMiddlewareTestApp(db)

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")

	// When: Making the request
	resp, err := app.Test(req)

	// Then: The request should be rejected with unauthorized status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Invalid token", response["error"])
}

func TestJWTMiddleware_ValidToken_UserNotFound(t *testing.T) {
	// Given: A request with a valid token but user doesn't exist in DB
	db := setupMiddlewareTestDB(t)
	app := setupMiddlewareTestApp(db)

	token, _ := security.GenerateToken(999, "nonexistent@test.com", "Non", "Existent")

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	// When: Making the request
	resp, err := app.Test(req)

	// Then: The request should be rejected with unauthorized status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "User not found", response["error"])
}

func TestJWTMiddleware_ValidToken_Success(t *testing.T) {
	// Given: A request with a valid token for an existing user
	db := setupMiddlewareTestDB(t)
	app := setupMiddlewareTestApp(db)

	user := models.User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@example.com",
		Password:  "hashedpassword",
	}
	db.Create(&user)

	token, _ := security.GenerateToken(user.ID, user.Email, user.FirstName, user.LastName)

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	// When: Making the request
	resp, err := app.Test(req)

	// Then: The request should succeed
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestJWTMiddleware_ValidToken_SetsUserInContext(t *testing.T) {
	// Given: A request with a valid token for an existing user
	db := setupMiddlewareTestDB(t)
	app := setupMiddlewareTestApp(db)

	user := models.User{
		FirstName: "Context",
		LastName:  "User",
		Email:     "context@example.com",
		Password:  "hashedpassword",
	}
	db.Create(&user)

	token, _ := security.GenerateToken(user.ID, user.Email, user.FirstName, user.LastName)

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	// When: Making the request
	resp, err := app.Test(req)

	// Then: The response should contain the user data from context
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "context@example.com", response["email"])
}

func TestJWTMiddleware_TokenWithoutBearerPrefix(t *testing.T) {
	// Given: A request with token without Bearer prefix
	db := setupMiddlewareTestDB(t)
	app := setupMiddlewareTestApp(db)

	user := models.User{
		FirstName: "No",
		LastName:  "Bearer",
		Email:     "nobearer@example.com",
		Password:  "hashedpassword",
	}
	db.Create(&user)

	token, _ := security.GenerateToken(user.ID, user.Email, user.FirstName, user.LastName)

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", token)

	// When: Making the request
	resp, err := app.Test(req)

	// Then: The request should still succeed (token is parsed without Bearer prefix)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestJWTMiddleware_EmptyAuthorizationHeader(t *testing.T) {
	// Given: A request with empty Authorization header value
	db := setupMiddlewareTestDB(t)
	app := setupMiddlewareTestApp(db)

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "")

	// When: Making the request
	resp, err := app.Test(req)

	// Then: The request should be rejected with unauthorized status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestJWTMiddleware_OnlyBearerPrefix(t *testing.T) {
	// Given: A request with only "Bearer " as authorization
	db := setupMiddlewareTestDB(t)
	app := setupMiddlewareTestApp(db)

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer ")

	// When: Making the request
	resp, err := app.Test(req)

	// Then: The request should be rejected with unauthorized status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestJWTMiddleware_MalformedToken(t *testing.T) {
	// Given: A request with a malformed token
	db := setupMiddlewareTestDB(t)
	app := setupMiddlewareTestApp(db)

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer not.a.valid.jwt.token")

	// When: Making the request
	resp, err := app.Test(req)

	// Then: The request should be rejected with unauthorized status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestUnauthorizedHandler(t *testing.T) {
	// Given: A fiber context and an error
	app := fiber.New()
	app.Get("/test", func(c *fiber.Ctx) error {
		return UnauthorizedHandler(c, fiber.ErrUnauthorized)
	})

	req := httptest.NewRequest("GET", "/test", nil)

	// When: Calling the unauthorized handler
	resp, err := app.Test(req)

	// Then: It should return unauthorized status with error message
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Unauthorized")
}

func TestJWTMiddleware_SetsUserIDInContext(t *testing.T) {
	// Given: A request with a valid token for an existing user
	db := setupMiddlewareTestDB(t)
	app := fiber.New()

	app.Get("/protected", JWTMiddleware(db), func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(uint)
		return c.JSON(fiber.Map{
			"userID": userID,
		})
	})

	user := models.User{
		FirstName: "ID",
		LastName:  "User",
		Email:     "id@example.com",
		Password:  "hashedpassword",
	}
	db.Create(&user)

	token, _ := security.GenerateToken(user.ID, user.Email, user.FirstName, user.LastName)

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	// When: Making the request
	resp, err := app.Test(req)

	// Then: The userID should be set in context
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, float64(user.ID), response["userID"])
}
