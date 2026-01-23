package dtos

type CreateCommentRequest struct {
	Content string `json:"content" validate:"required"`
	UserID  uint   `json:"user_id" validate:"required"`
	NewsID  uint   `json:"news_id" validate:"required"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" validate:"required"`
}

type CommentResponse struct {
	ID        uint   `json:"id"`
	Content   string `json:"content"`
	UserID    uint   `json:"user_id"`
	UserName  string `json:"user_name"`
	NewsID    uint   `json:"news_id"`
	NewsTitle string `json:"news_title"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
