package mappers

import (
	"testing"
	"time"

	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestToCommentResponse_BasicFields(t *testing.T) {
	// Given: A comment model with all fields populated
	now := time.Now()
	comment := &models.Comment{
		Model:   gorm.Model{ID: 1, CreatedAt: now, UpdatedAt: now},
		Content: "This is a test comment",
		UserID:  5,
		NewsID:  10,
		User:    models.User{FirstName: "John", LastName: "Doe"},
		News:    models.News{Title: "Test News"},
	}

	// When: Converting to response
	response := ToCommentResponse(comment)

	// Then: All fields should be correctly mapped
	assert.Equal(t, uint(1), response.ID)
	assert.Equal(t, "This is a test comment", response.Content)
	assert.Equal(t, uint(5), response.UserID)
	assert.Equal(t, "John Doe", response.UserName)
	assert.Equal(t, uint(10), response.NewsID)
	assert.Equal(t, "Test News", response.NewsTitle)
	assert.NotEmpty(t, response.CreatedAt)
	assert.NotEmpty(t, response.UpdatedAt)
}

func TestToCommentResponse_EmptyLastName(t *testing.T) {
	// Given: A comment model with user having only first name
	now := time.Now()
	comment := &models.Comment{
		Model:   gorm.Model{ID: 2, CreatedAt: now, UpdatedAt: now},
		Content: "Comment content",
		UserID:  1,
		NewsID:  1,
		User:    models.User{FirstName: "Jane", LastName: ""},
		News:    models.News{Title: "News Title"},
	}

	// When: Converting to response
	response := ToCommentResponse(comment)

	// Then: User name should handle empty last name
	assert.Equal(t, "Jane ", response.UserName)
}

func TestToCommentResponse_LongContent(t *testing.T) {
	// Given: A comment model with long content
	longContent := "This is a very long comment that contains multiple sentences. It should be preserved in its entirety when converting to a response. Lorem ipsum dolor sit amet."
	now := time.Now()
	comment := &models.Comment{
		Model:   gorm.Model{ID: 3, CreatedAt: now, UpdatedAt: now},
		Content: longContent,
		UserID:  1,
		NewsID:  1,
		User:    models.User{FirstName: "User"},
		News:    models.News{Title: "News"},
	}

	// When: Converting to response
	response := ToCommentResponse(comment)

	// Then: Long content should be preserved
	assert.Equal(t, longContent, response.Content)
}

func TestToCommentResponseList_EmptySlice(t *testing.T) {
	// Given: An empty slice of comments
	comments := []models.Comment{}

	// When: Converting to response list
	responses := ToCommentResponseList(comments)

	// Then: The result should be an empty slice
	assert.Empty(t, responses)
	assert.Len(t, responses, 0)
}

func TestToCommentResponseList_SingleComment(t *testing.T) {
	// Given: A slice with one comment
	now := time.Now()
	comments := []models.Comment{
		{
			Model:   gorm.Model{ID: 1, CreatedAt: now, UpdatedAt: now},
			Content: "Solo comment",
			UserID:  1,
			NewsID:  1,
			User:    models.User{FirstName: "Solo", LastName: "User"},
			News:    models.News{Title: "Solo News"},
		},
	}

	// When: Converting to response list
	responses := ToCommentResponseList(comments)

	// Then: One comment should be converted
	assert.Len(t, responses, 1)
	assert.Equal(t, "Solo comment", responses[0].Content)
	assert.Equal(t, "Solo User", responses[0].UserName)
}

func TestToCommentResponseList_MultipleComments(t *testing.T) {
	// Given: A slice with multiple comments
	now := time.Now()
	comments := []models.Comment{
		{Model: gorm.Model{ID: 1, CreatedAt: now, UpdatedAt: now}, Content: "Comment 1", User: models.User{FirstName: "User1"}, News: models.News{Title: "News1"}},
		{Model: gorm.Model{ID: 2, CreatedAt: now, UpdatedAt: now}, Content: "Comment 2", User: models.User{FirstName: "User2"}, News: models.News{Title: "News2"}},
		{Model: gorm.Model{ID: 3, CreatedAt: now, UpdatedAt: now}, Content: "Comment 3", User: models.User{FirstName: "User3"}, News: models.News{Title: "News3"}},
	}

	// When: Converting to response list
	responses := ToCommentResponseList(comments)

	// Then: All comments should be converted
	assert.Len(t, responses, 3)
	assert.Equal(t, "Comment 1", responses[0].Content)
	assert.Equal(t, "Comment 2", responses[1].Content)
	assert.Equal(t, "Comment 3", responses[2].Content)
}

func TestToCommentResponseList_PreservesOrder(t *testing.T) {
	// Given: A slice with comments in specific order
	now := time.Now()
	comments := []models.Comment{
		{Model: gorm.Model{ID: 3, CreatedAt: now, UpdatedAt: now}, Content: "Third", User: models.User{FirstName: "A"}, News: models.News{Title: "N"}},
		{Model: gorm.Model{ID: 1, CreatedAt: now, UpdatedAt: now}, Content: "First", User: models.User{FirstName: "B"}, News: models.News{Title: "N"}},
		{Model: gorm.Model{ID: 2, CreatedAt: now, UpdatedAt: now}, Content: "Second", User: models.User{FirstName: "C"}, News: models.News{Title: "N"}},
	}

	// When: Converting to response list
	responses := ToCommentResponseList(comments)

	// Then: The order should be preserved
	assert.Equal(t, "Third", responses[0].Content)
	assert.Equal(t, "First", responses[1].Content)
	assert.Equal(t, "Second", responses[2].Content)
}

func TestToCommentModel_BasicRequest(t *testing.T) {
	// Given: A create comment request
	req := dtos.CreateCommentRequest{
		Content: "New comment",
		UserID:  10,
		NewsID:  20,
	}

	// When: Converting to model
	comment := ToCommentModel(req)

	// Then: Fields should be set correctly
	assert.Equal(t, "New comment", comment.Content)
	assert.Equal(t, uint(10), comment.UserID)
	assert.Equal(t, uint(20), comment.NewsID)
}

func TestToCommentModel_SpecialCharactersInContent(t *testing.T) {
	// Given: A create comment request with special characters
	req := dtos.CreateCommentRequest{
		Content: "Special chars: <html> & symbols @#$%",
		UserID:  1,
		NewsID:  1,
	}

	// When: Converting to model
	comment := ToCommentModel(req)

	// Then: Special characters should be preserved
	assert.Equal(t, "Special chars: <html> & symbols @#$%", comment.Content)
}

func TestToCommentResponse_DateTimeFormat(t *testing.T) {
	// Given: A comment with specific timestamps
	createdAt := time.Date(2024, 7, 15, 10, 30, 45, 0, time.UTC)
	updatedAt := time.Date(2024, 7, 16, 14, 20, 30, 0, time.UTC)
	comment := &models.Comment{
		Model:   gorm.Model{ID: 1, CreatedAt: createdAt, UpdatedAt: updatedAt},
		Content: "Test",
		User:    models.User{FirstName: "Test"},
		News:    models.News{Title: "Test"},
	}

	// When: Converting to response
	response := ToCommentResponse(comment)

	// Then: Dates should be formatted correctly
	assert.Equal(t, "2024-07-15 10:30:45", response.CreatedAt)
	assert.Equal(t, "2024-07-16 14:20:30", response.UpdatedAt)
}
