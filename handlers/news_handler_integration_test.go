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

func setupNewsTestApp(db *gorm.DB) *fiber.App {
	app := fiber.New()
	newsHandler := NewNewsHandler(db)

	app.Get("/news", newsHandler.GetNews)
	app.Post("/news", newsHandler.CreateNews)
	app.Put("/news/:id", newsHandler.UpdateNews)
	app.Delete("/news/:id", newsHandler.DeleteNews)

	return app
}

func TestNewsHandler_GetNews_Empty(t *testing.T) {
	// Given: An empty database
	db := setupTestDB(t)
	app := setupNewsTestApp(db)

	req := httptest.NewRequest("GET", "/news", nil)

	// When: Making the get all news request
	resp, err := app.Test(req)

	// Then: The request should succeed with empty array
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var news []dtos.NewsResponse
	json.NewDecoder(resp.Body).Decode(&news)
	assert.Empty(t, news)
}

func TestNewsHandler_GetNews_WithNews(t *testing.T) {
	// Given: A database with users and news
	db := setupTestDB(t)
	app := setupNewsTestApp(db)

	user := models.User{FirstName: "John", LastName: "Doe", Email: "john@test.com", Password: "pass"}
	db.Create(&user)
	db.Create(&models.News{Title: "News 1", Description: "Description 1", AuthorID: user.ID, Date: "2024-01-01"})
	db.Create(&models.News{Title: "News 2", Description: "Description 2", AuthorID: user.ID, Date: "2024-01-02"})

	req := httptest.NewRequest("GET", "/news", nil)

	// When: Making the get all news request
	resp, err := app.Test(req)

	// Then: The request should return all news
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var news []dtos.NewsResponse
	json.NewDecoder(resp.Body).Decode(&news)
	assert.Len(t, news, 2)
}

func TestNewsHandler_CreateNews_Success(t *testing.T) {
	// Given: A valid create news request with existing author
	db := setupTestDB(t)
	app := setupNewsTestApp(db)

	user := models.User{FirstName: "Author", LastName: "User", Email: "author@test.com", Password: "pass"}
	db.Create(&user)

	reqBody := dtos.CreateNewsRequest{
		Title:       "New Article",
		Description: "Article content",
		AuthorId:    user.ID,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/news", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create news request
	resp, err := app.Test(req)

	// Then: The news should be created
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestNewsHandler_CreateNews_InvalidJSON(t *testing.T) {
	// Given: An invalid JSON request body
	db := setupTestDB(t)
	app := setupNewsTestApp(db)

	req := httptest.NewRequest("POST", "/news", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create news request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestNewsHandler_CreateNews_MissingAuthorId(t *testing.T) {
	// Given: A create news request without author ID
	db := setupTestDB(t)
	app := setupNewsTestApp(db)

	reqBody := dtos.CreateNewsRequest{
		Title:       "New Article",
		Description: "Content",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/news", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create news request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestNewsHandler_CreateNews_InvalidAuthorId(t *testing.T) {
	// Given: A create news request with non-existent author
	db := setupTestDB(t)
	app := setupNewsTestApp(db)

	reqBody := dtos.CreateNewsRequest{
		Title:       "New Article",
		Description: "Content",
		AuthorId:    999,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/news", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create news request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestNewsHandler_CreateNews_ReturnsCreatedNews(t *testing.T) {
	// Given: A valid create news request
	db := setupTestDB(t)
	app := setupNewsTestApp(db)

	user := models.User{FirstName: "Creator", LastName: "Test", Email: "creator@test.com", Password: "pass"}
	db.Create(&user)

	reqBody := dtos.CreateNewsRequest{
		Title:       "Created News",
		Description: "Created Description",
		AuthorId:    user.ID,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/news", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create news request
	resp, err := app.Test(req)

	// Then: The response should contain the created news
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var response dtos.NewsResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Created News", response.Title)
	assert.Equal(t, "Creator", response.Author)
}

func TestNewsHandler_UpdateNews_NotFound(t *testing.T) {
	// Given: An empty database
	db := setupTestDB(t)
	app := setupNewsTestApp(db)

	reqBody := dtos.CreateNewsRequest{Title: "Updated", Description: "Desc"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/news/999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update request for non-existent news
	resp, err := app.Test(req)

	// Then: The request should return not found
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestNewsHandler_UpdateNews_InvalidID(t *testing.T) {
	// Given: An invalid news ID
	db := setupTestDB(t)
	app := setupNewsTestApp(db)

	reqBody := dtos.CreateNewsRequest{Title: "Updated"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/news/invalid", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update request with invalid ID
	resp, err := app.Test(req)

	// Then: The request should return bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestNewsHandler_UpdateNews_InvalidJSON(t *testing.T) {
	// Given: An existing news and invalid JSON
	db := setupTestDB(t)
	app := setupNewsTestApp(db)

	user := models.User{FirstName: "Author", Email: "a@test.com", Password: "pass"}
	db.Create(&user)
	news := models.News{Title: "Original", AuthorID: user.ID, Date: "2024-01-01"}
	db.Create(&news)

	req := httptest.NewRequest("PUT", fmt.Sprintf("/news/%d", news.ID), bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update request with invalid JSON
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestNewsHandler_DeleteNews_Success(t *testing.T) {
	// Given: An existing news article
	db := setupTestDB(t)
	app := setupNewsTestApp(db)

	user := models.User{FirstName: "Author", Email: "a@test.com", Password: "pass"}
	db.Create(&user)
	news := models.News{Title: "To Delete", AuthorID: user.ID, Date: "2024-01-01"}
	db.Create(&news)

	req := httptest.NewRequest("DELETE", fmt.Sprintf("/news/%d", news.ID), nil)

	// When: Making the delete news request
	resp, err := app.Test(req)

	// Then: The news should be deleted
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
}

func TestNewsHandler_DeleteNews_NotFound(t *testing.T) {
	// Given: An empty database
	db := setupTestDB(t)
	app := setupNewsTestApp(db)

	req := httptest.NewRequest("DELETE", "/news/999", nil)

	// When: Making the delete request for non-existent news
	resp, err := app.Test(req)

	// Then: The request should return not found
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestNewsHandler_DeleteNews_InvalidID(t *testing.T) {
	// Given: An invalid news ID
	db := setupTestDB(t)
	app := setupNewsTestApp(db)

	req := httptest.NewRequest("DELETE", "/news/invalid", nil)

	// When: Making the delete request with invalid ID
	resp, err := app.Test(req)

	// Then: The request should return bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestNewsHandler_DeleteNews_ActuallyDeletes(t *testing.T) {
	// Given: An existing news article
	db := setupTestDB(t)
	app := setupNewsTestApp(db)

	user := models.User{FirstName: "Author", Email: "a@test.com", Password: "pass"}
	db.Create(&user)
	news := models.News{Title: "To Delete", AuthorID: user.ID, Date: "2024-01-01"}
	db.Create(&news)
	newsID := news.ID

	req := httptest.NewRequest("DELETE", fmt.Sprintf("/news/%d", newsID), nil)

	// When: Making the delete news request
	app.Test(req)

	// Then: The news should no longer exist in the database
	var count int64
	db.Model(&models.News{}).Where("id = ?", newsID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestNewNewsHandler(t *testing.T) {
	// Given: A database connection
	db := setupTestDB(t)

	// When: Creating a new news handler
	handler := NewNewsHandler(db)

	// Then: The handler should be created
	assert.NotNil(t, handler)
}
