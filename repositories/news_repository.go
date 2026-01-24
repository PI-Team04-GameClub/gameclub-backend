package repositories

import (
	"context"

	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"gorm.io/gorm"
)

const newsWhereIDEquals = "id = ?"

type NewsRepository interface {
	FindAll(ctx context.Context) ([]models.News, error)
	FindByID(ctx context.Context, id int) (*models.News, error)
	Create(ctx context.Context, news *models.News) error
	Update(ctx context.Context, news *models.News) error
	Delete(ctx context.Context, id int) error
}

type newsRepository struct {
	db *gorm.DB
}

func NewNewsRepository(db *gorm.DB) NewsRepository {
	return &newsRepository{db: db}
}

func (r *newsRepository) FindAll(ctx context.Context) ([]models.News, error) {
	return gorm.G[models.News](r.db).Preload("Author", nil).Find(ctx)
}

func (r *newsRepository) FindByID(ctx context.Context, id int) (*models.News, error) {
	news, err := gorm.G[models.News](r.db).Preload("Author", nil).Where(newsWhereIDEquals, id).First(ctx)
	if err != nil {
		return nil, err
	}
	return &news, nil
}

func (r *newsRepository) Create(ctx context.Context, news *models.News) error {
	return gorm.G[models.News](r.db).Create(ctx, news)
}

func (r *newsRepository) Update(ctx context.Context, news *models.News) error {
	_, err := gorm.G[models.News](r.db).Where(newsWhereIDEquals, news.ID).Updates(ctx, *news)
	return err
}

func (r *newsRepository) Delete(ctx context.Context, id int) error {
	_, err := gorm.G[models.News](r.db).Where(newsWhereIDEquals, id).Delete(ctx)
	return err
}
