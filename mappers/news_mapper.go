package mappers

import (
	"time"

	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
)

func ToNewsResponse(news *models.News, authorName string) dtos.NewsResponse {
	return dtos.NewsResponse{
		ID:          int(news.ID),
		Title:       news.Title,
		Description: news.Description,
		Author:      authorName,
		Date:        news.Date,
	}
}

func ToNewsModel(req dtos.CreateNewsRequest, authorID uint) models.News {
	return models.News{
		Title:       req.Title,
		Description: req.Description,
		AuthorID:    authorID,
		Date:        time.Now().Format("2006-01-02"),
	}
}

func ToNewsResponseList(newsList []models.News) []dtos.NewsResponse {
	responses := make([]dtos.NewsResponse, len(newsList))
	for i, news := range newsList {
		responses[i] = ToNewsResponse(&news, news.Author.FirstName)
	}
	return responses
}
