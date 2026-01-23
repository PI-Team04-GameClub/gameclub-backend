package mocks

import (
	"context"

	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/stretchr/testify/mock"
)

type MockNewsRepository struct {
	mock.Mock
}

func (m *MockNewsRepository) mockMethodError(methodName string, args ...interface{}) error {
	return m.MethodCalled(methodName, args...).Error(0)
}

func (m *MockNewsRepository) FindAll(ctx context.Context) ([]models.News, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.News), args.Error(1)
}

func (m *MockNewsRepository) FindByID(ctx context.Context, id int) (*models.News, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.News), args.Error(1)
}

func (m *MockNewsRepository) Create(ctx context.Context, news *models.News) error {
	return m.mockMethodError("Create", ctx, news)
}

func (m *MockNewsRepository) Update(ctx context.Context, news *models.News) error {
	return m.mockMethodError("Update", ctx, news)
}

func (m *MockNewsRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
