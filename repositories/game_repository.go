package repositories

import (
	"context"

	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"gorm.io/gorm"
)

type GameRepository interface {
	FindAll(ctx context.Context) ([]models.Game, error)
	FindByID(ctx context.Context, id string) (*models.Game, error)
	Create(ctx context.Context, game *models.Game) error
	Update(ctx context.Context, game *models.Game) error
	Delete(ctx context.Context, id uint) error
}

type gameRepository struct {
	db *gorm.DB
}

func NewGameRepository(db *gorm.DB) GameRepository {
	return &gameRepository{db: db}
}

func (r *gameRepository) FindAll(ctx context.Context) ([]models.Game, error) {
	return gorm.G[models.Game](r.db).Find(ctx)
}

func (r *gameRepository) FindByID(ctx context.Context, id string) (*models.Game, error) {
	game, err := gorm.G[models.Game](r.db).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func (r *gameRepository) Create(ctx context.Context, game *models.Game) error {
	return gorm.G[models.Game](r.db).Create(ctx, game)
}

func (r *gameRepository) Update(ctx context.Context, game *models.Game) error {
	_, err := gorm.G[models.Game](r.db).Where("id = ?", game.ID).Updates(ctx, *game)
	return err
}

func (r *gameRepository) Delete(ctx context.Context, id uint) error {
	_, err := gorm.G[models.Game](r.db).Where("id = ?", id).Delete(ctx)
	return err
}
