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

func setupCommentTestApp(db *gorm.DB) *fiber.App {
	app := fiber.New()
	commentHandler := NewCommentHandler(db)

	app.Post("/comments", commentHandler.CreateComment)
	app.Put("/comments/:id", commentHandler.UpdateComment)
	app.Get("/users/:id/comments", commentHandler.GetCommentsByUserID)
	app.Get("/news/:id/comments", commentHandler.GetCommentsByNewsID)

	return app
}

func createTestUserAndNews(db *gorm.DB) (models.User, models.News) {
	user := models.User{FirstName: "Test", LastName: "User", Email: "test@test.com", Password: "pass"}
	db.Create(&user)

	news := models.News{Title: "Test News", Description: "Test Description", AuthorID: user.ID, Date: "2024-01-01"}
	db.Create(&news)

	return user, news
}

func TestCommentHandler_CreateComment_Success(t *testing.T) {
	// Given: A valid create comment request with existing user and news
	db := setupTestDB(t)
	db.AutoMigrate(&models.Comment{})
	app := setupCommentTestApp(db)

	user, news := createTestUserAndNews(db)

	reqBody := dtos.CreateCommentRequest{
		Content: "This is a test comment",
		UserID:  user.ID,
		NewsID:  news.ID,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/comments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create comment request
	resp, err := app.Test(req)

	// Then: The comment should be created
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var response dtos.CommentResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "This is a test comment", response.Content)
	assert.Equal(t, user.ID, response.UserID)
	assert.Equal(t, news.ID, response.NewsID)
}

func TestCommentHandler_CreateComment_InvalidJSON(t *testing.T) {
	// Given: An invalid JSON request body
	db := setupTestDB(t)
	db.AutoMigrate(&models.Comment{})
	app := setupCommentTestApp(db)

	req := httptest.NewRequest("POST", "/comments", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create comment request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestCommentHandler_CreateComment_MissingContent(t *testing.T) {
	// Given: A create comment request without content
	db := setupTestDB(t)
	db.AutoMigrate(&models.Comment{})
	app := setupCommentTestApp(db)

	user, news := createTestUserAndNews(db)

	reqBody := dtos.CreateCommentRequest{
		UserID: user.ID,
		NewsID: news.ID,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/comments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create comment request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestCommentHandler_CreateComment_UserNotFound(t *testing.T) {
	// Given: A create comment request with non-existent user
	db := setupTestDB(t)
	db.AutoMigrate(&models.Comment{})
	app := setupCommentTestApp(db)

	_, news := createTestUserAndNews(db)

	reqBody := dtos.CreateCommentRequest{
		Content: "Test comment",
		UserID:  999,
		NewsID:  news.ID,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/comments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create comment request
	resp, err := app.Test(req)

	// Then: The request should fail with not found
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestCommentHandler_CreateComment_NewsNotFound(t *testing.T) {
	// Given: A create comment request with non-existent news
	db := setupTestDB(t)
	db.AutoMigrate(&models.Comment{})
	app := setupCommentTestApp(db)

	user, _ := createTestUserAndNews(db)

	reqBody := dtos.CreateCommentRequest{
		Content: "Test comment",
		UserID:  user.ID,
		NewsID:  999,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/comments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create comment request
	resp, err := app.Test(req)

	// Then: The request should fail with not found
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestCommentHandler_UpdateComment_Success(t *testing.T) {
	// Given: An existing comment
	db := setupTestDB(t)
	db.AutoMigrate(&models.Comment{})
	app := setupCommentTestApp(db)

	user, news := createTestUserAndNews(db)
	comment := models.Comment{Content: "Original content", UserID: user.ID, NewsID: news.ID}
	db.Create(&comment)

	reqBody := dtos.UpdateCommentRequest{
		Content: "Updated content",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", fmt.Sprintf("/comments/%d", comment.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update comment request
	resp, err := app.Test(req)

	// Then: The comment should be updated
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response dtos.CommentResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Updated content", response.Content)
}

func TestCommentHandler_UpdateComment_NotFound(t *testing.T) {
	// Given: An empty database
	db := setupTestDB(t)
	db.AutoMigrate(&models.Comment{})
	app := setupCommentTestApp(db)

	reqBody := dtos.UpdateCommentRequest{Content: "Updated"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/comments/999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update request for non-existent comment
	resp, err := app.Test(req)

	// Then: The request should return not found
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestCommentHandler_UpdateComment_InvalidID(t *testing.T) {
	// Given: An invalid comment ID
	db := setupTestDB(t)
	db.AutoMigrate(&models.Comment{})
	app := setupCommentTestApp(db)

	reqBody := dtos.UpdateCommentRequest{Content: "Updated"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/comments/invalid", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update request with invalid ID
	resp, err := app.Test(req)

	// Then: The request should return bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestCommentHandler_UpdateComment_MissingContent(t *testing.T) {
	// Given: An existing comment and empty content
	db := setupTestDB(t)
	db.AutoMigrate(&models.Comment{})
	app := setupCommentTestApp(db)

	user, news := createTestUserAndNews(db)
	comment := models.Comment{Content: "Original content", UserID: user.ID, NewsID: news.ID}
	db.Create(&comment)

	reqBody := dtos.UpdateCommentRequest{Content: ""}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", fmt.Sprintf("/comments/%d", comment.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update comment request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestCommentHandler_GetCommentsByUserID_Success(t *testing.T) {
	// Given: A user with multiple comments
	db := setupTestDB(t)
	db.AutoMigrate(&models.Comment{})
	app := setupCommentTestApp(db)

	user, news := createTestUserAndNews(db)
	db.Create(&models.Comment{Content: "Comment 1", UserID: user.ID, NewsID: news.ID})
	db.Create(&models.Comment{Content: "Comment 2", UserID: user.ID, NewsID: news.ID})

	req := httptest.NewRequest("GET", fmt.Sprintf("/users/%d/comments", user.ID), nil)

	// When: Making the get comments by user ID request
	resp, err := app.Test(req)

	// Then: The request should return all user's comments
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var comments []dtos.CommentResponse
	json.NewDecoder(resp.Body).Decode(&comments)
	assert.Len(t, comments, 2)
}

func TestCommentHandler_GetCommentsByUserID_Empty(t *testing.T) {
	// Given: A user with no comments
	db := setupTestDB(t)
	db.AutoMigrate(&models.Comment{})
	app := setupCommentTestApp(db)

	user := models.User{FirstName: "No", LastName: "Comments", Email: "nocomments@test.com", Password: "pass"}
	db.Create(&user)

	req := httptest.NewRequest("GET", fmt.Sprintf("/users/%d/comments", user.ID), nil)

	// When: Making the get comments by user ID request
	resp, err := app.Test(req)

	// Then: The request should return empty array
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var comments []dtos.CommentResponse
	json.NewDecoder(resp.Body).Decode(&comments)
	assert.Empty(t, comments)
}

func TestCommentHandler_GetCommentsByUserID_UserNotFound(t *testing.T) {
	// Given: An empty database
	db := setupTestDB(t)
	db.AutoMigrate(&models.Comment{})
	app := setupCommentTestApp(db)

	req := httptest.NewRequest("GET", "/users/999/comments", nil)

	// When: Making the get comments by user ID request
	resp, err := app.Test(req)

	// Then: The request should return not found
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestCommentHandler_GetCommentsByUserID_InvalidID(t *testing.T) {
	// Given: An invalid user ID
	db := setupTestDB(t)
	db.AutoMigrate(&models.Comment{})
	app := setupCommentTestApp(db)

	req := httptest.NewRequest("GET", "/users/invalid/comments", nil)

	// When: Making the get comments by user ID request
	resp, err := app.Test(req)

	// Then: The request should return bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestCommentHandler_GetCommentsByNewsID_Success(t *testing.T) {
	// Given: A news article with multiple comments
	db := setupTestDB(t)
	db.AutoMigrate(&models.Comment{})
	app := setupCommentTestApp(db)

	user, news := createTestUserAndNews(db)
	user2 := models.User{FirstName: "Another", LastName: "User", Email: "another@test.com", Password: "pass"}
	db.Create(&user2)

	db.Create(&models.Comment{Content: "Comment 1", UserID: user.ID, NewsID: news.ID})
	db.Create(&models.Comment{Content: "Comment 2", UserID: user2.ID, NewsID: news.ID})

	req := httptest.NewRequest("GET", fmt.Sprintf("/news/%d/comments", news.ID), nil)

	// When: Making the get comments by news ID request
	resp, err := app.Test(req)

	// Then: The request should return all news' comments
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var comments []dtos.CommentResponse
	json.NewDecoder(resp.Body).Decode(&comments)
	assert.Len(t, comments, 2)
}

func TestCommentHandler_GetCommentsByNewsID_Empty(t *testing.T) {
	// Given: A news article with no comments
	db := setupTestDB(t)
	db.AutoMigrate(&models.Comment{})
	app := setupCommentTestApp(db)

	user := models.User{FirstName: "Author", LastName: "User", Email: "author@test.com", Password: "pass"}
	db.Create(&user)
	news := models.News{Title: "No Comments News", Description: "Desc", AuthorID: user.ID, Date: "2024-01-01"}
	db.Create(&news)

	req := httptest.NewRequest("GET", fmt.Sprintf("/news/%d/comments", news.ID), nil)

	// When: Making the get comments by news ID request
	resp, err := app.Test(req)

	// Then: The request should return empty array
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var comments []dtos.CommentResponse
	json.NewDecoder(resp.Body).Decode(&comments)
	assert.Empty(t, comments)
}

func TestCommentHandler_GetCommentsByNewsID_NewsNotFound(t *testing.T) {
	// Given: An empty database
	db := setupTestDB(t)
	db.AutoMigrate(&models.Comment{})
	app := setupCommentTestApp(db)

	req := httptest.NewRequest("GET", "/news/999/comments", nil)

	// When: Making the get comments by news ID request
	resp, err := app.Test(req)

	// Then: The request should return not found
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestCommentHandler_GetCommentsByNewsID_InvalidID(t *testing.T) {
	// Given: An invalid news ID
	db := setupTestDB(t)
	db.AutoMigrate(&models.Comment{})
	app := setupCommentTestApp(db)

	req := httptest.NewRequest("GET", "/news/invalid/comments", nil)

	// When: Making the get comments by news ID request
	resp, err := app.Test(req)

	// Then: The request should return bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestNewCommentHandler(t *testing.T) {
	// Given: A database connection
	db := setupTestDB(t)

	// When: Creating a new comment handler
	handler := NewCommentHandler(db)

	// Then: The handler should be created
	assert.NotNil(t, handler)
}

func TestCommentHandler_CreateComment_ReturnsUserAndNewsInfo(t *testing.T) {
	// Given: A valid create comment request
	db := setupTestDB(t)
	db.AutoMigrate(&models.Comment{})
	app := setupCommentTestApp(db)

	user, news := createTestUserAndNews(db)

	reqBody := dtos.CreateCommentRequest{
		Content: "Test comment",
		UserID:  user.ID,
		NewsID:  news.ID,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/comments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create comment request
	resp, err := app.Test(req)

	// Then: The response should include user and news info
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var response dtos.CommentResponse
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Test User", response.UserName)
	assert.Equal(t, "Test News", response.NewsTitle)
}
