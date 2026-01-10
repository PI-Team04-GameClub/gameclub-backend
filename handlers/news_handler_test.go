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

func TestNewsHandler_GetNews_Success_Unit(t *testing.T) {
	// Given: News items exist in the repository
	mockNewsRepo := new(mocks.MockNewsRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewNewsHandlerWithRepo(mockNewsRepo, mockUserRepo)

	app := fiber.New()
	app.Get("/news", handler.GetNews)

	newsList := []models.News{
		{Model: gorm.Model{ID: 1}, Title: "News 1", Description: "Description 1"},
		{Model: gorm.Model{ID: 2}, Title: "News 2", Description: "Description 2"},
	}
	mockNewsRepo.On("FindAll", mock.Anything).Return(newsList, nil)

	req := httptest.NewRequest("GET", "/news", nil)

	// When: Making the get all news request
	resp, err := app.Test(req)

	// Then: The request should succeed with news list
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockNewsRepo.AssertExpectations(t)
}

func TestNewsHandler_GetNews_Empty_Unit(t *testing.T) {
	// Given: No news items exist in the repository
	mockNewsRepo := new(mocks.MockNewsRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewNewsHandlerWithRepo(mockNewsRepo, mockUserRepo)

	app := fiber.New()
	app.Get("/news", handler.GetNews)

	mockNewsRepo.On("FindAll", mock.Anything).Return([]models.News{}, nil)

	req := httptest.NewRequest("GET", "/news", nil)

	// When: Making the get all news request
	resp, err := app.Test(req)

	// Then: The request should succeed with empty list
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockNewsRepo.AssertExpectations(t)
}

func TestNewsHandler_GetNews_DatabaseError_Unit(t *testing.T) {
	// Given: A database error occurs
	mockNewsRepo := new(mocks.MockNewsRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewNewsHandlerWithRepo(mockNewsRepo, mockUserRepo)

	app := fiber.New()
	app.Get("/news", handler.GetNews)

	mockNewsRepo.On("FindAll", mock.Anything).Return(nil, errors.New("database error"))

	req := httptest.NewRequest("GET", "/news", nil)

	// When: Making the get all news request
	resp, err := app.Test(req)

	// Then: The request should fail with internal server error
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockNewsRepo.AssertExpectations(t)
}

func TestNewsHandler_CreateNews_Success_Unit(t *testing.T) {
	// Given: A valid create news request
	mockNewsRepo := new(mocks.MockNewsRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewNewsHandlerWithRepo(mockNewsRepo, mockUserRepo)

	app := fiber.New()
	app.Post("/news", handler.CreateNews)

	author := &models.User{Model: gorm.Model{ID: 1}, FirstName: "John", LastName: "Doe"}
	mockUserRepo.On("FindByID", mock.Anything, uint(1)).Return(author, nil)
	mockNewsRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.News")).Return(nil)
	mockNewsRepo.On("FindByID", mock.Anything, 0).Return(&models.News{
		Model:       gorm.Model{ID: 1},
		Title:       "New News",
		Description: "News Description",
		Author:      *author,
	}, nil)

	reqBody := dtos.CreateNewsRequest{
		Title:       "New News",
		Description: "News Description",
		AuthorId:    1,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/news", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create news request
	resp, err := app.Test(req)

	// Then: The request should succeed with created status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	mockNewsRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestNewsHandler_CreateNews_InvalidJSON_Unit(t *testing.T) {
	// Given: An invalid JSON request body
	mockNewsRepo := new(mocks.MockNewsRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewNewsHandlerWithRepo(mockNewsRepo, mockUserRepo)

	app := fiber.New()
	app.Post("/news", handler.CreateNews)

	req := httptest.NewRequest("POST", "/news", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create news request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestNewsHandler_CreateNews_MissingAuthorID_Unit(t *testing.T) {
	// Given: A request without author ID
	mockNewsRepo := new(mocks.MockNewsRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewNewsHandlerWithRepo(mockNewsRepo, mockUserRepo)

	app := fiber.New()
	app.Post("/news", handler.CreateNews)

	reqBody := dtos.CreateNewsRequest{
		Title:       "New News",
		Description: "News Description",
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

func TestNewsHandler_CreateNews_AuthorNotFound_Unit(t *testing.T) {
	// Given: The referenced author does not exist
	mockNewsRepo := new(mocks.MockNewsRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewNewsHandlerWithRepo(mockNewsRepo, mockUserRepo)

	app := fiber.New()
	app.Post("/news", handler.CreateNews)

	mockUserRepo.On("FindByID", mock.Anything, uint(999)).Return(nil, gorm.ErrRecordNotFound)

	reqBody := dtos.CreateNewsRequest{
		Title:       "New News",
		Description: "News Description",
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
	mockUserRepo.AssertExpectations(t)
}

func TestNewsHandler_CreateNews_DatabaseError_Unit(t *testing.T) {
	// Given: A database error occurs when creating news
	mockNewsRepo := new(mocks.MockNewsRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewNewsHandlerWithRepo(mockNewsRepo, mockUserRepo)

	app := fiber.New()
	app.Post("/news", handler.CreateNews)

	author := &models.User{Model: gorm.Model{ID: 1}, FirstName: "John", LastName: "Doe"}
	mockUserRepo.On("FindByID", mock.Anything, uint(1)).Return(author, nil)
	mockNewsRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.News")).Return(errors.New("database error"))

	reqBody := dtos.CreateNewsRequest{
		Title:       "New News",
		Description: "News Description",
		AuthorId:    1,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/news", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the create news request
	resp, err := app.Test(req)

	// Then: The request should fail with internal server error
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockNewsRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestNewsHandler_UpdateNews_Success_Unit(t *testing.T) {
	// Given: A news item exists and valid update request
	mockNewsRepo := new(mocks.MockNewsRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewNewsHandlerWithRepo(mockNewsRepo, mockUserRepo)

	app := fiber.New()
	app.Put("/news/:id", handler.UpdateNews)

	author := models.User{Model: gorm.Model{ID: 1}, FirstName: "John", LastName: "Doe"}
	existingNews := &models.News{Model: gorm.Model{ID: 1}, Title: "Old Title", Description: "Old Desc", Author: author}
	mockNewsRepo.On("FindByID", mock.Anything, 1).Return(existingNews, nil)
	mockNewsRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.News")).Return(nil)

	reqBody := dtos.CreateNewsRequest{
		Title:       "New Title",
		Description: "New Description",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/news/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update news request
	resp, err := app.Test(req)

	// Then: The request should succeed
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockNewsRepo.AssertExpectations(t)
}

func TestNewsHandler_UpdateNews_NotFound_Unit(t *testing.T) {
	// Given: No news item exists with the specified ID
	mockNewsRepo := new(mocks.MockNewsRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewNewsHandlerWithRepo(mockNewsRepo, mockUserRepo)

	app := fiber.New()
	app.Put("/news/:id", handler.UpdateNews)

	mockNewsRepo.On("FindByID", mock.Anything, 999).Return(nil, gorm.ErrRecordNotFound)

	reqBody := dtos.CreateNewsRequest{
		Title:       "New Title",
		Description: "New Description",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/news/999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update news request
	resp, err := app.Test(req)

	// Then: The request should fail with not found status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockNewsRepo.AssertExpectations(t)
}

func TestNewsHandler_UpdateNews_InvalidID_Unit(t *testing.T) {
	// Given: An invalid news ID
	mockNewsRepo := new(mocks.MockNewsRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewNewsHandlerWithRepo(mockNewsRepo, mockUserRepo)

	app := fiber.New()
	app.Put("/news/:id", handler.UpdateNews)

	reqBody := dtos.CreateNewsRequest{
		Title:       "New Title",
		Description: "New Description",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/news/invalid", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update news request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestNewsHandler_UpdateNews_InvalidJSON_Unit(t *testing.T) {
	// Given: A news item exists but request body is invalid JSON
	mockNewsRepo := new(mocks.MockNewsRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewNewsHandlerWithRepo(mockNewsRepo, mockUserRepo)

	app := fiber.New()
	app.Put("/news/:id", handler.UpdateNews)

	author := models.User{Model: gorm.Model{ID: 1}, FirstName: "John", LastName: "Doe"}
	existingNews := &models.News{Model: gorm.Model{ID: 1}, Title: "Old Title", Description: "Old Desc", Author: author}
	mockNewsRepo.On("FindByID", mock.Anything, 1).Return(existingNews, nil)

	req := httptest.NewRequest("PUT", "/news/1", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update news request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	mockNewsRepo.AssertExpectations(t)
}

func TestNewsHandler_UpdateNews_DatabaseError_Unit(t *testing.T) {
	// Given: A news item exists but database error occurs during update
	mockNewsRepo := new(mocks.MockNewsRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewNewsHandlerWithRepo(mockNewsRepo, mockUserRepo)

	app := fiber.New()
	app.Put("/news/:id", handler.UpdateNews)

	author := models.User{Model: gorm.Model{ID: 1}, FirstName: "John", LastName: "Doe"}
	existingNews := &models.News{Model: gorm.Model{ID: 1}, Title: "Old Title", Description: "Old Desc", Author: author}
	mockNewsRepo.On("FindByID", mock.Anything, 1).Return(existingNews, nil)
	mockNewsRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.News")).Return(errors.New("database error"))

	reqBody := dtos.CreateNewsRequest{
		Title:       "New Title",
		Description: "New Description",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/news/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// When: Making the update news request
	resp, err := app.Test(req)

	// Then: The request should fail with internal server error
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockNewsRepo.AssertExpectations(t)
}

func TestNewsHandler_DeleteNews_Success_Unit(t *testing.T) {
	// Given: A news item exists with the specified ID
	mockNewsRepo := new(mocks.MockNewsRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewNewsHandlerWithRepo(mockNewsRepo, mockUserRepo)

	app := fiber.New()
	app.Delete("/news/:id", handler.DeleteNews)

	existingNews := &models.News{Model: gorm.Model{ID: 1}, Title: "News to Delete"}
	mockNewsRepo.On("FindByID", mock.Anything, 1).Return(existingNews, nil)
	mockNewsRepo.On("Delete", mock.Anything, 1).Return(nil)

	req := httptest.NewRequest("DELETE", "/news/1", nil)

	// When: Making the delete news request
	resp, err := app.Test(req)

	// Then: The request should succeed with no content status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	mockNewsRepo.AssertExpectations(t)
}

func TestNewsHandler_DeleteNews_NotFound_Unit(t *testing.T) {
	// Given: No news item exists with the specified ID
	mockNewsRepo := new(mocks.MockNewsRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewNewsHandlerWithRepo(mockNewsRepo, mockUserRepo)

	app := fiber.New()
	app.Delete("/news/:id", handler.DeleteNews)

	mockNewsRepo.On("FindByID", mock.Anything, 999).Return(nil, gorm.ErrRecordNotFound)

	req := httptest.NewRequest("DELETE", "/news/999", nil)

	// When: Making the delete news request
	resp, err := app.Test(req)

	// Then: The request should fail with not found status
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockNewsRepo.AssertExpectations(t)
}

func TestNewsHandler_DeleteNews_InvalidID_Unit(t *testing.T) {
	// Given: An invalid news ID
	mockNewsRepo := new(mocks.MockNewsRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewNewsHandlerWithRepo(mockNewsRepo, mockUserRepo)

	app := fiber.New()
	app.Delete("/news/:id", handler.DeleteNews)

	req := httptest.NewRequest("DELETE", "/news/invalid", nil)

	// When: Making the delete news request
	resp, err := app.Test(req)

	// Then: The request should fail with bad request
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestNewsHandler_DeleteNews_DatabaseError_Unit(t *testing.T) {
	// Given: A news item exists but database error occurs during delete
	mockNewsRepo := new(mocks.MockNewsRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	handler := NewNewsHandlerWithRepo(mockNewsRepo, mockUserRepo)

	app := fiber.New()
	app.Delete("/news/:id", handler.DeleteNews)

	existingNews := &models.News{Model: gorm.Model{ID: 1}, Title: "News to Delete"}
	mockNewsRepo.On("FindByID", mock.Anything, 1).Return(existingNews, nil)
	mockNewsRepo.On("Delete", mock.Anything, 1).Return(errors.New("database error"))

	req := httptest.NewRequest("DELETE", "/news/1", nil)

	// When: Making the delete news request
	resp, err := app.Test(req)

	// Then: The request should fail with internal server error
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockNewsRepo.AssertExpectations(t)
}
