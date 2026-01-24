package mocks

import (
	"context"

	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/stretchr/testify/mock"
)

type MockFriendRequestRepository struct {
	mock.Mock
}

func (m *MockFriendRequestRepository) mockMethodError(methodName string, args ...interface{}) error {
	return m.MethodCalled(methodName, args...).Error(0)
}

func getResultOrNil[T any](args mock.Arguments) (T, error) {
	if args.Get(0) == nil {
		var zero T
		return zero, args.Error(1)
	}
	return args.Get(0).(T), args.Error(1)
}

func (m *MockFriendRequestRepository) Create(ctx context.Context, friendRequest *models.FriendRequest) error {
	return m.mockMethodError("Create", ctx, friendRequest)
}

func (m *MockFriendRequestRepository) FindByID(ctx context.Context, id uint) (*models.FriendRequest, error) {
	return getResultOrNil[*models.FriendRequest](m.Called(ctx, id))
}

func (m *MockFriendRequestRepository) Update(ctx context.Context, friendRequest *models.FriendRequest) error {
	return m.mockMethodError("Update", ctx, friendRequest)
}

func (m *MockFriendRequestRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockFriendRequestRepository) FindBySenderID(ctx context.Context, senderID uint) ([]models.FriendRequest, error) {
	return getResultOrNil[[]models.FriendRequest](m.Called(ctx, senderID))
}

func (m *MockFriendRequestRepository) FindByReceiverID(ctx context.Context, receiverID uint) ([]models.FriendRequest, error) {
	return getResultOrNil[[]models.FriendRequest](m.Called(ctx, receiverID))
}

func (m *MockFriendRequestRepository) FindPendingByReceiverID(ctx context.Context, receiverID uint) ([]models.FriendRequest, error) {
	return getResultOrNil[[]models.FriendRequest](m.Called(ctx, receiverID))
}

func (m *MockFriendRequestRepository) FindFriendsByUserID(ctx context.Context, userID uint) ([]models.User, error) {
	return getResultOrNil[[]models.User](m.Called(ctx, userID))
}

func (m *MockFriendRequestRepository) FindByUsers(ctx context.Context, userID1, userID2 uint) (*models.FriendRequest, error) {
	return getResultOrNil[*models.FriendRequest](m.Called(ctx, userID1, userID2))
}
