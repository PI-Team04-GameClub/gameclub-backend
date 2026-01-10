package mocks

import (
	"context"

	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/stretchr/testify/mock"
)

type MockGameRepository struct {
	mock.Mock
}

func (m *MockGameRepository) FindAll(ctx context.Context) ([]models.Game, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Game), args.Error(1)
}

func (m *MockGameRepository) FindByID(ctx context.Context, id string) (*models.Game, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Game), args.Error(1)
}

func (m *MockGameRepository) Create(ctx context.Context, game *models.Game) error {
	args := m.Called(ctx, game)
	return args.Error(0)
}

func (m *MockGameRepository) Update(ctx context.Context, game *models.Game) error {
	args := m.Called(ctx, game)
	return args.Error(0)
}

func (m *MockGameRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
