package repositories

import (
	"context"

	"github.com/PI-Team04-GameClub/gameclub-backend/models"

	"gorm.io/gorm"
)

const (
	preloadUser = "User"
	preloadNews = "News"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *models.Comment) error
	FindByID(ctx context.Context, id uint) (*models.Comment, error)
	Update(ctx context.Context, comment *models.Comment) error
	Delete(ctx context.Context, id uint) error
	FindByUserID(ctx context.Context, userID uint) ([]models.Comment, error)
	FindByNewsID(ctx context.Context, newsID uint) ([]models.Comment, error)
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(ctx context.Context, comment *models.Comment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

func (r *commentRepository) FindByID(ctx context.Context, id uint) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.WithContext(ctx).Preload(preloadUser).Preload(preloadNews).First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) Update(ctx context.Context, comment *models.Comment) error {
	return r.db.WithContext(ctx).Save(comment).Error
}

func (r *commentRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Comment{}, id).Error
}

func (r *commentRepository) FindByUserID(ctx context.Context, userID uint) ([]models.Comment, error) {
	var comments []models.Comment
	err := r.db.WithContext(ctx).Preload(preloadUser).Preload(preloadNews).Where("user_id = ?", userID).Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *commentRepository) FindByNewsID(ctx context.Context, newsID uint) ([]models.Comment, error) {
	var comments []models.Comment
	err := r.db.WithContext(ctx).Preload(preloadUser).Preload(preloadNews).Where("news_id = ?", newsID).Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}
