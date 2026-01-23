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

func TestCommentHandler_CreateComment_Success_Unit(t *testing.T) {
	// Given: A valid create comment request
	mockCommentRepo := new(mocks.MockCommentRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNewsRepo := new(mocks.MockNewsRepository)
	handler := NewCommentHandlerWithRepo(mockCommentRepo, mockUserRepo, mockNewsRepo)

	app := fiber.New()
	app.Post("/comments", handler.CreateComment)

	user := &models.User{Model: gorm.Model{ID: 1}, FirstName: "John", LastName: "Doe"}
	news := &models.News{Model: gorm.Model{ID: 1}, Title: "Test News"}
	mockUserRepo.On("FindByID", mock.Anything, uint(1)).Return(user, nil)
	mockNewsRepo.On("FindByID", mock.Anything, 1).Return(news, nil)
	mockCommentRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Comment")).Return(nil)
	mockCommentRepo.On("FindByID", mock.Anything, uint(0)).Return(&models.Comment{
		Model:   gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Content: "Test comment",
		UserID:  1,
		NewsID:  1,
		User:    *user,
		News:    *news,
	}, nil)

	reqBody := dtos.CreateCommentRequest{
		Content: "Test comment",
		UserID:  1,
		NewsID:  1,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/comments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create comment request
	resp, err := app.Test(req)

	// Then: The request should succeed with created status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	mockCommentRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockNewsRepo.AssertExpectations(t)
}

func TestCommentHandler_CreateComment_InvalidJSON_Unit(t *testing.T) {
	// Given: An invalid JSON request body
	mockCommentRepo := new(mocks.MockCommentRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNewsRepo := new(mocks.MockNewsRepository)
	handler := NewCommentHandlerWithRepo(mockCommentRepo, mockUserRepo, mockNewsRepo)

	app := fiber.New()
	app.Post("/comments", handler.CreateComment)

	req := httptest.NewRequest("POST", "/comments", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create comment request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestCommentHandler_CreateComment_MissingContent_Unit(t *testing.T) {
	// Given: A request without content
	mockCommentRepo := new(mocks.MockCommentRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNewsRepo := new(mocks.MockNewsRepository)
	handler := NewCommentHandlerWithRepo(mockCommentRepo, mockUserRepo, mockNewsRepo)

	app := fiber.New()
	app.Post("/comments", handler.CreateComment)

	reqBody := dtos.CreateCommentRequest{
		UserID: 1,
		NewsID: 1,
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

func TestCommentHandler_CreateComment_MissingUserID_Unit(t *testing.T) {
	// Given: A request without user ID
	mockCommentRepo := new(mocks.MockCommentRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNewsRepo := new(mocks.MockNewsRepository)
	handler := NewCommentHandlerWithRepo(mockCommentRepo, mockUserRepo, mockNewsRepo)

	app := fiber.New()
	app.Post("/comments", handler.CreateComment)

	reqBody := dtos.CreateCommentRequest{
		Content: "Test comment",
		NewsID:  1,
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

func TestCommentHandler_CreateComment_MissingNewsID_Unit(t *testing.T) {
	// Given: A request without news ID
	mockCommentRepo := new(mocks.MockCommentRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNewsRepo := new(mocks.MockNewsRepository)
	handler := NewCommentHandlerWithRepo(mockCommentRepo, mockUserRepo, mockNewsRepo)

	app := fiber.New()
	app.Post("/comments", handler.CreateComment)

	reqBody := dtos.CreateCommentRequest{
		Content: "Test comment",
		UserID:  1,
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

func TestCommentHandler_CreateComment_UserNotFound_Unit(t *testing.T) {
	// Given: The referenced user does not exist
	mockCommentRepo := new(mocks.MockCommentRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNewsRepo := new(mocks.MockNewsRepository)
	handler := NewCommentHandlerWithRepo(mockCommentRepo, mockUserRepo, mockNewsRepo)

	app := fiber.New()
	app.Post("/comments", handler.CreateComment)

	mockUserRepo.On("FindByID", mock.Anything, uint(999)).Return(nil, gorm.ErrRecordNotFound)

	reqBody := dtos.CreateCommentRequest{
		Content: "Test comment",
		UserID:  999,
		NewsID:  1,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/comments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create comment request
	resp, err := app.Test(req)

	// Then: The request should fail with not found
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

func TestCommentHandler_CreateComment_NewsNotFound_Unit(t *testing.T) {
	// Given: The referenced news does not exist
	mockCommentRepo := new(mocks.MockCommentRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNewsRepo := new(mocks.MockNewsRepository)
	handler := NewCommentHandlerWithRepo(mockCommentRepo, mockUserRepo, mockNewsRepo)

	app := fiber.New()
	app.Post("/comments", handler.CreateComment)

	user := &models.User{Model: gorm.Model{ID: 1}, FirstName: "John"}
	mockUserRepo.On("FindByID", mock.Anything, uint(1)).Return(user, nil)
	mockNewsRepo.On("FindByID", mock.Anything, 999).Return(nil, gorm.ErrRecordNotFound)

	reqBody := dtos.CreateCommentRequest{
		Content: "Test comment",
		UserID:  1,
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
	mockUserRepo.AssertExpectations(t)
	mockNewsRepo.AssertExpectations(t)
}

func TestCommentHandler_CreateComment_DatabaseError_Unit(t *testing.T) {
	// Given: A database error occurs when creating comment
	mockCommentRepo := new(mocks.MockCommentRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNewsRepo := new(mocks.MockNewsRepository)
	handler := NewCommentHandlerWithRepo(mockCommentRepo, mockUserRepo, mockNewsRepo)

	app := fiber.New()
	app.Post("/comments", handler.CreateComment)

	user := &models.User{Model: gorm.Model{ID: 1}, FirstName: "John"}
	news := &models.News{Model: gorm.Model{ID: 1}, Title: "Test News"}
	mockUserRepo.On("FindByID", mock.Anything, uint(1)).Return(user, nil)
	mockNewsRepo.On("FindByID", mock.Anything, 1).Return(news, nil)
	mockCommentRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Comment")).Return(errors.New("database error"))

	reqBody := dtos.CreateCommentRequest{
		Content: "Test comment",
		UserID:  1,
		NewsID:  1,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/comments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create comment request
	resp, err := app.Test(req)

	// Then: The request should fail with internal server error
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockCommentRepo.AssertExpectations(t)
}

func TestCommentHandler_UpdateComment_Success_Unit(t *testing.T) {
	// Given: A comment exists and valid update request
	mockCommentRepo := new(mocks.MockCommentRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNewsRepo := new(mocks.MockNewsRepository)
	handler := NewCommentHandlerWithRepo(mockCommentRepo, mockUserRepo, mockNewsRepo)

	app := fiber.New()
	app.Put("/comments/:id", handler.UpdateComment)

	user := models.User{Model: gorm.Model{ID: 1}, FirstName: "John", LastName: "Doe"}
	news := models.News{Model: gorm.Model{ID: 1}, Title: "Test News"}
	existingComment := &models.Comment{
		Model:   gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Content: "Old content",
		UserID:  1,
		NewsID:  1,
		User:    user,
		News:    news,
	}
	mockCommentRepo.On("FindByID", mock.Anything, uint(1)).Return(existingComment, nil).Times(2)
	mockCommentRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Comment")).Return(nil)

	reqBody := dtos.UpdateCommentRequest{
		Content: "Updated content",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/comments/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update comment request
	resp, err := app.Test(req)

	// Then: The request should succeed
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockCommentRepo.AssertExpectations(t)
}

func TestCommentHandler_UpdateComment_NotFound_Unit(t *testing.T) {
	// Given: No comment exists with the specified ID
	mockCommentRepo := new(mocks.MockCommentRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNewsRepo := new(mocks.MockNewsRepository)
	handler := NewCommentHandlerWithRepo(mockCommentRepo, mockUserRepo, mockNewsRepo)

	app := fiber.New()
	app.Put("/comments/:id", handler.UpdateComment)

	mockCommentRepo.On("FindByID", mock.Anything, uint(999)).Return(nil, gorm.ErrRecordNotFound)

	reqBody := dtos.UpdateCommentRequest{
		Content: "Updated content",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/comments/999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update comment request
	resp, err := app.Test(req)

	// Then: The request should fail with not found status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockCommentRepo.AssertExpectations(t)
}

func TestCommentHandler_UpdateComment_InvalidID_Unit(t *testing.T) {
	// Given: An invalid comment ID
	mockCommentRepo := new(mocks.MockCommentRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNewsRepo := new(mocks.MockNewsRepository)
	handler := NewCommentHandlerWithRepo(mockCommentRepo, mockUserRepo, mockNewsRepo)

	app := fiber.New()
	app.Put("/comments/:id", handler.UpdateComment)

	reqBody := dtos.UpdateCommentRequest{
		Content: "Updated content",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/comments/invalid", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update comment request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestCommentHandler_UpdateComment_MissingContent_Unit(t *testing.T) {
	// Given: A request without content
	mockCommentRepo := new(mocks.MockCommentRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNewsRepo := new(mocks.MockNewsRepository)
	handler := NewCommentHandlerWithRepo(mockCommentRepo, mockUserRepo, mockNewsRepo)

	app := fiber.New()
	app.Put("/comments/:id", handler.UpdateComment)

	user := models.User{Model: gorm.Model{ID: 1}, FirstName: "John"}
	news := models.News{Model: gorm.Model{ID: 1}, Title: "Test News"}
	existingComment := &models.Comment{
		Model:   gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Content: "Old content",
		User:    user,
		News:    news,
	}
	mockCommentRepo.On("FindByID", mock.Anything, uint(1)).Return(existingComment, nil)

	reqBody := dtos.UpdateCommentRequest{
		Content: "",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/comments/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update comment request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestCommentHandler_UpdateComment_DatabaseError_Unit(t *testing.T) {
	// Given: A database error occurs during update
	mockCommentRepo := new(mocks.MockCommentRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNewsRepo := new(mocks.MockNewsRepository)
	handler := NewCommentHandlerWithRepo(mockCommentRepo, mockUserRepo, mockNewsRepo)

	app := fiber.New()
	app.Put("/comments/:id", handler.UpdateComment)

	user := models.User{Model: gorm.Model{ID: 1}, FirstName: "John"}
	news := models.News{Model: gorm.Model{ID: 1}, Title: "Test News"}
	existingComment := &models.Comment{
		Model:   gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		Content: "Old content",
		User:    user,
		News:    news,
	}
	mockCommentRepo.On("FindByID", mock.Anything, uint(1)).Return(existingComment, nil)
	mockCommentRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Comment")).Return(errors.New("database error"))

	reqBody := dtos.UpdateCommentRequest{
		Content: "Updated content",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/comments/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update comment request
	resp, err := app.Test(req)

	// Then: The request should fail with internal server error
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockCommentRepo.AssertExpectations(t)
}

func TestCommentHandler_GetCommentsByUserID_Success_Unit(t *testing.T) {
	// Given: A user exists with comments
	mockCommentRepo := new(mocks.MockCommentRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNewsRepo := new(mocks.MockNewsRepository)
	handler := NewCommentHandlerWithRepo(mockCommentRepo, mockUserRepo, mockNewsRepo)

	app := fiber.New()
	app.Get("/users/:id/comments", handler.GetCommentsByUserID)

	user := &models.User{Model: gorm.Model{ID: 1}, FirstName: "John"}
	mockUserRepo.On("FindByID", mock.Anything, uint(1)).Return(user, nil)

	comments := []models.Comment{
		{Model: gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Content: "Comment 1", User: *user, News: models.News{Title: "News 1"}},
		{Model: gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Content: "Comment 2", User: *user, News: models.News{Title: "News 2"}},
	}
	mockCommentRepo.On("FindByUserID", mock.Anything, uint(1)).Return(comments, nil)

	req := httptest.NewRequest("GET", "/users/1/comments", nil)

	// When: Making the get comments by user ID request
	resp, err := app.Test(req)

	// Then: The request should succeed
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
	mockCommentRepo.AssertExpectations(t)
}

func TestCommentHandler_GetCommentsByUserID_UserNotFound_Unit(t *testing.T) {
	// Given: The user does not exist
	mockCommentRepo := new(mocks.MockCommentRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNewsRepo := new(mocks.MockNewsRepository)
	handler := NewCommentHandlerWithRepo(mockCommentRepo, mockUserRepo, mockNewsRepo)

	app := fiber.New()
	app.Get("/users/:id/comments", handler.GetCommentsByUserID)

	mockUserRepo.On("FindByID", mock.Anything, uint(999)).Return(nil, gorm.ErrRecordNotFound)

	req := httptest.NewRequest("GET", "/users/999/comments", nil)

	// When: Making the get comments by user ID request
	resp, err := app.Test(req)

	// Then: The request should fail with not found
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockUserRepo.AssertExpectations(t)
}

func TestCommentHandler_GetCommentsByUserID_InvalidID_Unit(t *testing.T) {
	// Given: An invalid user ID
	mockCommentRepo := new(mocks.MockCommentRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNewsRepo := new(mocks.MockNewsRepository)
	handler := NewCommentHandlerWithRepo(mockCommentRepo, mockUserRepo, mockNewsRepo)

	app := fiber.New()
	app.Get("/users/:id/comments", handler.GetCommentsByUserID)

	req := httptest.NewRequest("GET", "/users/invalid/comments", nil)

	// When: Making the get comments by user ID request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestCommentHandler_GetCommentsByNewsID_Success_Unit(t *testing.T) {
	// Given: A news exists with comments
	mockCommentRepo := new(mocks.MockCommentRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNewsRepo := new(mocks.MockNewsRepository)
	handler := NewCommentHandlerWithRepo(mockCommentRepo, mockUserRepo, mockNewsRepo)

	app := fiber.New()
	app.Get("/news/:id/comments", handler.GetCommentsByNewsID)

	news := &models.News{Model: gorm.Model{ID: 1}, Title: "Test News"}
	mockNewsRepo.On("FindByID", mock.Anything, 1).Return(news, nil)

	comments := []models.Comment{
		{Model: gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Content: "Comment 1", User: models.User{FirstName: "User1"}, News: *news},
		{Model: gorm.Model{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Content: "Comment 2", User: models.User{FirstName: "User2"}, News: *news},
	}
	mockCommentRepo.On("FindByNewsID", mock.Anything, uint(1)).Return(comments, nil)

	req := httptest.NewRequest("GET", "/news/1/comments", nil)

	// When: Making the get comments by news ID request
	resp, err := app.Test(req)

	// Then: The request should succeed
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockNewsRepo.AssertExpectations(t)
	mockCommentRepo.AssertExpectations(t)
}

func TestCommentHandler_GetCommentsByNewsID_NewsNotFound_Unit(t *testing.T) {
	// Given: The news does not exist
	mockCommentRepo := new(mocks.MockCommentRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNewsRepo := new(mocks.MockNewsRepository)
	handler := NewCommentHandlerWithRepo(mockCommentRepo, mockUserRepo, mockNewsRepo)

	app := fiber.New()
	app.Get("/news/:id/comments", handler.GetCommentsByNewsID)

	mockNewsRepo.On("FindByID", mock.Anything, 999).Return(nil, gorm.ErrRecordNotFound)

	req := httptest.NewRequest("GET", "/news/999/comments", nil)

	// When: Making the get comments by news ID request
	resp, err := app.Test(req)

	// Then: The request should fail with not found
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockNewsRepo.AssertExpectations(t)
}

func TestCommentHandler_GetCommentsByNewsID_InvalidID_Unit(t *testing.T) {
	// Given: An invalid news ID
	mockCommentRepo := new(mocks.MockCommentRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNewsRepo := new(mocks.MockNewsRepository)
	handler := NewCommentHandlerWithRepo(mockCommentRepo, mockUserRepo, mockNewsRepo)

	app := fiber.New()
	app.Get("/news/:id/comments", handler.GetCommentsByNewsID)

	req := httptest.NewRequest("GET", "/news/invalid/comments", nil)

	// When: Making the get comments by news ID request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestNewCommentHandler_Unit(t *testing.T) {
	// Given: A database connection
	mockCommentRepo := new(mocks.MockCommentRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockNewsRepo := new(mocks.MockNewsRepository)

	// When: Creating a new comment handler with repos
	handler := NewCommentHandlerWithRepo(mockCommentRepo, mockUserRepo, mockNewsRepo)

	// Then: The handler should be created
	assert.NotNil(t, handler)
}
