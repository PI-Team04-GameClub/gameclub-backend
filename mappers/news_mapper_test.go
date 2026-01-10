package mappers

import (
	"testing"
	"time"

	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestToNewsResponse_BasicFields(t *testing.T) {
	// Given: A news model with all fields populated
	news := &models.News{
		Model:       gorm.Model{ID: 1},
		Title:       "Breaking News",
		Description: "This is a test news article",
		AuthorID:    5,
		Date:        "2024-07-15",
	}
	authorName := "John Doe"

	// When: Converting to response
	response := ToNewsResponse(news, authorName)

	// Then: All fields should be correctly mapped
	assert.Equal(t, 1, response.ID)
	assert.Equal(t, "Breaking News", response.Title)
	assert.Equal(t, "This is a test news article", response.Description)
	assert.Equal(t, "John Doe", response.Author)
	assert.Equal(t, "2024-07-15", response.Date)
}

func TestToNewsResponse_EmptyDescription(t *testing.T) {
	// Given: A news model with empty description
	news := &models.News{
		Model:       gorm.Model{ID: 2},
		Title:       "Title Only",
		Description: "",
		Date:        "2024-08-01",
	}

	// When: Converting to response
	response := ToNewsResponse(news, "Author")

	// Then: Empty description should be preserved
	assert.Equal(t, "", response.Description)
}

func TestToNewsResponse_LongContent(t *testing.T) {
	// Given: A news model with long content
	longDescription := "This is a very long description that contains multiple sentences. It should be preserved in its entirety when converting to a response. Lorem ipsum dolor sit amet."
	news := &models.News{
		Model:       gorm.Model{ID: 3},
		Title:       "Long Article",
		Description: longDescription,
		Date:        "2024-09-01",
	}

	// When: Converting to response
	response := ToNewsResponse(news, "Writer")

	// Then: Long content should be preserved
	assert.Equal(t, longDescription, response.Description)
}

func TestToNewsModel_BasicRequest(t *testing.T) {
	// Given: A create news request and author ID
	req := dtos.CreateNewsRequest{
		Title:       "New Article",
		Description: "Article description",
		AuthorId:    10,
	}
	authorID := uint(10)

	// When: Converting to model
	news := ToNewsModel(req, authorID)

	// Then: Fields should be set correctly
	assert.Equal(t, "New Article", news.Title)
	assert.Equal(t, "Article description", news.Description)
	assert.Equal(t, uint(10), news.AuthorID)
}

func TestToNewsModel_AutoGeneratesDate(t *testing.T) {
	// Given: A create news request
	req := dtos.CreateNewsRequest{
		Title:       "Dated Article",
		Description: "Test",
		AuthorId:    1,
	}

	// When: Converting to model
	news := ToNewsModel(req, 1)

	// Then: Date should be auto-generated in YYYY-MM-DD format
	assert.NotEmpty(t, news.Date)
	today := time.Now().Format("2006-01-02")
	assert.Equal(t, today, news.Date)
}

func TestToNewsModel_SpecialCharactersInTitle(t *testing.T) {
	// Given: A create news request with special characters
	req := dtos.CreateNewsRequest{
		Title:       "Special: News & Updates (2024)",
		Description: "Content with <html> tags & symbols",
		AuthorId:    5,
	}

	// When: Converting to model
	news := ToNewsModel(req, 5)

	// Then: Special characters should be preserved
	assert.Equal(t, "Special: News & Updates (2024)", news.Title)
	assert.Equal(t, "Content with <html> tags & symbols", news.Description)
}

func TestToNewsResponseList_EmptySlice(t *testing.T) {
	// Given: An empty slice of news
	newsList := []models.News{}

	// When: Converting to response list
	responses := ToNewsResponseList(newsList)

	// Then: The result should be an empty slice
	assert.Empty(t, responses)
	assert.Len(t, responses, 0)
}

func TestToNewsResponseList_SingleNews(t *testing.T) {
	// Given: A slice with one news item with author
	newsList := []models.News{
		{
			Model:       gorm.Model{ID: 1},
			Title:       "Solo News",
			Description: "Only one",
			Date:        "2024-01-01",
			Author:      models.User{FirstName: "Solo Author"},
		},
	}

	// When: Converting to response list
	responses := ToNewsResponseList(newsList)

	// Then: One news should be converted with correct author
	assert.Len(t, responses, 1)
	assert.Equal(t, "Solo News", responses[0].Title)
	assert.Equal(t, "Solo Author", responses[0].Author)
}

func TestToNewsResponseList_MultipleNews(t *testing.T) {
	// Given: A slice with multiple news items
	newsList := []models.News{
		{Model: gorm.Model{ID: 1}, Title: "News 1", Date: "2024-01-01", Author: models.User{FirstName: "Author1"}},
		{Model: gorm.Model{ID: 2}, Title: "News 2", Date: "2024-01-02", Author: models.User{FirstName: "Author2"}},
		{Model: gorm.Model{ID: 3}, Title: "News 3", Date: "2024-01-03", Author: models.User{FirstName: "Author3"}},
	}

	// When: Converting to response list
	responses := ToNewsResponseList(newsList)

	// Then: All news should be converted
	assert.Len(t, responses, 3)
	assert.Equal(t, "News 1", responses[0].Title)
	assert.Equal(t, "News 2", responses[1].Title)
	assert.Equal(t, "News 3", responses[2].Title)
}

func TestToNewsResponseList_PreservesOrder(t *testing.T) {
	// Given: A slice with news in specific order
	newsList := []models.News{
		{Model: gorm.Model{ID: 3}, Title: "Third", Date: "2024-03-01", Author: models.User{FirstName: "A"}},
		{Model: gorm.Model{ID: 1}, Title: "First", Date: "2024-01-01", Author: models.User{FirstName: "B"}},
		{Model: gorm.Model{ID: 2}, Title: "Second", Date: "2024-02-01", Author: models.User{FirstName: "C"}},
	}

	// When: Converting to response list
	responses := ToNewsResponseList(newsList)

	// Then: The order should be preserved
	assert.Equal(t, "Third", responses[0].Title)
	assert.Equal(t, "First", responses[1].Title)
	assert.Equal(t, "Second", responses[2].Title)
}

func TestToNewsResponse_IDIsInt(t *testing.T) {
	// Given: A news model with uint ID
	news := &models.News{
		Model: gorm.Model{ID: 999},
		Title: "ID Test",
		Date:  "2024-01-01",
	}

	// When: Converting to response
	response := ToNewsResponse(news, "Author")

	// Then: ID should be converted to int
	assert.Equal(t, 999, response.ID)
}

func TestToNewsModel_DifferentAuthorIDs(t *testing.T) {
	// Given: A request with one author ID and a different parameter
	req := dtos.CreateNewsRequest{
		Title:    "Test",
		AuthorId: 100,
	}
	authorIDParam := uint(200)

	// When: Converting to model
	news := ToNewsModel(req, authorIDParam)

	// Then: The authorID parameter should be used
	assert.Equal(t, uint(200), news.AuthorID)
}

func TestToNewsResponseList_UsesAuthorFirstName(t *testing.T) {
	// Given: A slice with news where author has first and last name
	newsList := []models.News{
		{
			Model: gorm.Model{ID: 1},
			Title: "Author Test",
			Date:  "2024-01-01",
			Author: models.User{
				FirstName: "John",
				LastName:  "Doe",
			},
		},
	}

	// When: Converting to response list
	responses := ToNewsResponseList(newsList)

	// Then: Only the first name should be used as author
	assert.Equal(t, "John", responses[0].Author)
}
