package mocks

import (
	"context"

	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/stretchr/testify/mock"
)

type MockFriendRequestRepository struct {
	mock.Mock
}

func (m *MockFriendRequestRepository) Create(ctx context.Context, friendRequest *models.FriendRequest) error {
	args := m.Called(ctx, friendRequest)
	return args.Error(0)
}

func (m *MockFriendRequestRepository) FindByID(ctx context.Context, id uint) (*models.FriendRequest, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FriendRequest), args.Error(1)
}

func (m *MockFriendRequestRepository) Update(ctx context.Context, friendRequest *models.FriendRequest) error {
	args := m.Called(ctx, friendRequest)
	return args.Error(0)
}

func (m *MockFriendRequestRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockFriendRequestRepository) FindBySenderID(ctx context.Context, senderID uint) ([]models.FriendRequest, error) {
	args := m.Called(ctx, senderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.FriendRequest), args.Error(1)
}

func (m *MockFriendRequestRepository) FindByReceiverID(ctx context.Context, receiverID uint) ([]models.FriendRequest, error) {
	args := m.Called(ctx, receiverID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.FriendRequest), args.Error(1)
}

func (m *MockFriendRequestRepository) FindPendingByReceiverID(ctx context.Context, receiverID uint) ([]models.FriendRequest, error) {
	args := m.Called(ctx, receiverID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.FriendRequest), args.Error(1)
}

func (m *MockFriendRequestRepository) FindFriendsByUserID(ctx context.Context, userID uint) ([]models.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockFriendRequestRepository) FindByUsers(ctx context.Context, userID1, userID2 uint) (*models.FriendRequest, error) {
	args := m.Called(ctx, userID1, userID2)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FriendRequest), args.Error(1)
}
