package repositories

import (
	"context"

	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"gorm.io/gorm"
)

const tournamentWhereIDEquals = "id = ?"

type TournamentRepository interface {
	FindAll(ctx context.Context) ([]models.Tournament, error)
	FindByID(ctx context.Context, id int) (*models.Tournament, error)
	Create(ctx context.Context, tournament *models.Tournament) error
	Update(ctx context.Context, tournament *models.Tournament) error
	Delete(ctx context.Context, id int) error
}

type tournamentRepository struct {
	db *gorm.DB
}

func NewTournamentRepository(db *gorm.DB) TournamentRepository {
	return &tournamentRepository{db: db}
}

func (r *tournamentRepository) FindAll(ctx context.Context) ([]models.Tournament, error) {
	return gorm.G[models.Tournament](r.db).Preload("Game", nil).Find(ctx)
}

func (r *tournamentRepository) FindByID(ctx context.Context, id int) (*models.Tournament, error) {
	tournament, err := gorm.G[models.Tournament](r.db).Preload("Game", nil).Where(tournamentWhereIDEquals, id).First(ctx)
	if err != nil {
		return nil, err
	}
	return &tournament, nil
}

func (r *tournamentRepository) Create(ctx context.Context, tournament *models.Tournament) error {
	return gorm.G[models.Tournament](r.db).Create(ctx, tournament)
}

func (r *tournamentRepository) Update(ctx context.Context, tournament *models.Tournament) error {
	_, err := gorm.G[models.Tournament](r.db).Where(tournamentWhereIDEquals, tournament.ID).Updates(ctx, *tournament)
	return err
}

func (r *tournamentRepository) Delete(ctx context.Context, id int) error {
	_, err := gorm.G[models.Tournament](r.db).Where(tournamentWhereIDEquals, id).Delete(ctx)
	return err
}
