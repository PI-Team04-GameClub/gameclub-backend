package repositories

import (
	"context"

	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"gorm.io/gorm"
)

type TeamRepository interface {
	FindAll(ctx context.Context) ([]models.Team, error)
	FindByID(ctx context.Context, id string) (*models.Team, error)
	FindByIDWithMembers(ctx context.Context, id string) (*models.Team, error)
	Create(ctx context.Context, team *models.Team) error
	Update(ctx context.Context, team *models.Team) error
	Delete(ctx context.Context, id uint) error
	AddMember(ctx context.Context, team *models.Team, user *models.User) error
}

type teamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepository{db: db}
}

func (r *teamRepository) FindAll(ctx context.Context) ([]models.Team, error) {
	return gorm.G[models.Team](r.db).Find(ctx)
}

func (r *teamRepository) FindByID(ctx context.Context, id string) (*models.Team, error) {
	team, err := gorm.G[models.Team](r.db).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *teamRepository) Create(ctx context.Context, team *models.Team) error {
	return gorm.G[models.Team](r.db).Create(ctx, team)
}

func (r *teamRepository) Update(ctx context.Context, team *models.Team) error {
	_, err := gorm.G[models.Team](r.db).Where("id = ?", team.ID).Updates(ctx, *team)
	return err
}

func (r *teamRepository) Delete(ctx context.Context, id uint) error {
	_, err := gorm.G[models.Team](r.db).Where("id = ?", id).Delete(ctx)
	return err
}

func (r *teamRepository) FindByIDWithMembers(ctx context.Context, id string) (*models.Team, error) {
	team, err := gorm.G[models.Team](r.db).Preload("Users", nil).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *teamRepository) AddMember(ctx context.Context, team *models.Team, user *models.User) error {
	return r.db.WithContext(ctx).Model(team).Association("Users").Append(user)
}
