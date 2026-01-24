package mocks

import (
	"context"

	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/stretchr/testify/mock"
)

type MockTournamentRepository struct {
	mock.Mock
}

func (m *MockTournamentRepository) mockMethodError(methodName string, args ...interface{}) error {
	return m.MethodCalled(methodName, args...).Error(0)
}

func (m *MockTournamentRepository) FindAll(ctx context.Context) ([]models.Tournament, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Tournament), args.Error(1)
}

func (m *MockTournamentRepository) FindByID(ctx context.Context, id int) (*models.Tournament, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tournament), args.Error(1)
}

func (m *MockTournamentRepository) Create(ctx context.Context, tournament *models.Tournament) error {
	return m.mockMethodError("Create", ctx, tournament)
}

func (m *MockTournamentRepository) Update(ctx context.Context, tournament *models.Tournament) error {
	return m.mockMethodError("Update", ctx, tournament)
}

func (m *MockTournamentRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
