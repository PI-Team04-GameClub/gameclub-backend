package mappers

import (
	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
)

func ToCommentResponse(comment *models.Comment) dtos.CommentResponse {
	return dtos.CommentResponse{
		ID:        comment.ID,
		Content:   comment.Content,
		UserID:    comment.UserID,
		UserName:  comment.User.FirstName + " " + comment.User.LastName,
		NewsID:    comment.NewsID,
		NewsTitle: comment.News.Title,
		CreatedAt: comment.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: comment.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToCommentResponseList(comments []models.Comment) []dtos.CommentResponse {
	responses := make([]dtos.CommentResponse, len(comments))
	for i, comment := range comments {
		responses[i] = ToCommentResponse(&comment)
	}
	return responses
}

func ToCommentModel(req dtos.CreateCommentRequest) models.Comment {
	return models.Comment{
		Content: req.Content,
		UserID:  req.UserID,
		NewsID:  req.NewsID,
	}
}
