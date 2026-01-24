package mocks

import (
	"context"

	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/stretchr/testify/mock"
)

type MockTeamRepository struct {
	mock.Mock
}

func (m *MockTeamRepository) mockMethodError(methodName string, args ...interface{}) error {
	return m.MethodCalled(methodName, args...).Error(0)
}

func (m *MockTeamRepository) FindAll(ctx context.Context) ([]models.Team, error) {
	return getResultOrNil[[]models.Team](m.Called(ctx))
}

func (m *MockTeamRepository) FindByID(ctx context.Context, id string) (*models.Team, error) {
	return getResultOrNil[*models.Team](m.Called(ctx, id))
}

func (m *MockTeamRepository) Create(ctx context.Context, team *models.Team) error {
	return m.mockMethodError("Create", ctx, team)
}

func (m *MockTeamRepository) Update(ctx context.Context, team *models.Team) error {
	return m.mockMethodError("Update", ctx, team)
}

func (m *MockTeamRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTeamRepository) FindByIDWithMembers(ctx context.Context, id string) (*models.Team, error) {
	return getResultOrNil[*models.Team](m.Called(ctx, id))
}

func (m *MockTeamRepository) AddMember(ctx context.Context, team *models.Team, user *models.User) error {
	args := m.Called(ctx, team, user)
	return args.Error(0)
}
