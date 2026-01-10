package repositories

import (
	"context"

	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id uint) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	FindAll(ctx context.Context) ([]models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := gorm.G[models.User](r.db).Where("email = ?", email).First(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	user, err := gorm.G[models.User](r.db).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return gorm.G[models.User](r.db).Create(ctx, user)
}

func (r *userRepository) FindAll(ctx context.Context) ([]models.User, error) {
	return gorm.G[models.User](r.db).Find(ctx)
}
