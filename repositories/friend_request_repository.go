package repositories

import (
	"context"

	"github.com/PI-Team04-GameClub/gameclub-backend/models"

	"gorm.io/gorm"
)

const (
	preloadSender   = "Sender"
	preloadReceiver = "Receiver"
)

type FriendRequestRepository interface {
	Create(ctx context.Context, friendRequest *models.FriendRequest) error
	FindByID(ctx context.Context, id uint) (*models.FriendRequest, error)
	Update(ctx context.Context, friendRequest *models.FriendRequest) error
	Delete(ctx context.Context, id uint) error
	FindBySenderID(ctx context.Context, senderID uint) ([]models.FriendRequest, error)
	FindByReceiverID(ctx context.Context, receiverID uint) ([]models.FriendRequest, error)
	FindPendingByReceiverID(ctx context.Context, receiverID uint) ([]models.FriendRequest, error)
	FindFriendsByUserID(ctx context.Context, userID uint) ([]models.User, error)
	FindByUsers(ctx context.Context, userID1, userID2 uint) (*models.FriendRequest, error)
}

type friendRequestRepository struct {
	db *gorm.DB
}

func NewFriendRequestRepository(db *gorm.DB) FriendRequestRepository {
	return &friendRequestRepository{db: db}
}

func (r *friendRequestRepository) Create(ctx context.Context, friendRequest *models.FriendRequest) error {
	return r.db.WithContext(ctx).Create(friendRequest).Error
}

func (r *friendRequestRepository) FindByID(ctx context.Context, id uint) (*models.FriendRequest, error) {
	var friendRequest models.FriendRequest
	err := r.db.WithContext(ctx).Preload(preloadSender).Preload(preloadReceiver).First(&friendRequest, id).Error
	if err != nil {
		return nil, err
	}
	return &friendRequest, nil
}

func (r *friendRequestRepository) Update(ctx context.Context, friendRequest *models.FriendRequest) error {
	return r.db.WithContext(ctx).Save(friendRequest).Error
}

func (r *friendRequestRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.FriendRequest{}, id).Error
}

func (r *friendRequestRepository) FindBySenderID(ctx context.Context, senderID uint) ([]models.FriendRequest, error) {
	var friendRequests []models.FriendRequest
	err := r.db.WithContext(ctx).Preload(preloadSender).Preload(preloadReceiver).Where("sender_id = ?", senderID).Find(&friendRequests).Error
	if err != nil {
		return nil, err
	}
	return friendRequests, nil
}

func (r *friendRequestRepository) FindByReceiverID(ctx context.Context, receiverID uint) ([]models.FriendRequest, error) {
	var friendRequests []models.FriendRequest
	err := r.db.WithContext(ctx).Preload(preloadSender).Preload(preloadReceiver).Where("receiver_id = ?", receiverID).Find(&friendRequests).Error
	if err != nil {
		return nil, err
	}
	return friendRequests, nil
}

func (r *friendRequestRepository) FindPendingByReceiverID(ctx context.Context, receiverID uint) ([]models.FriendRequest, error) {
	var friendRequests []models.FriendRequest
	err := r.db.WithContext(ctx).Preload(preloadSender).Preload(preloadReceiver).
		Where("receiver_id = ? AND status = ?", receiverID, models.StatusPending).
		Find(&friendRequests).Error
	if err != nil {
		return nil, err
	}
	return friendRequests, nil
}

func (r *friendRequestRepository) FindFriendsByUserID(ctx context.Context, userID uint) ([]models.User, error) {
	var friends []models.User

	err := r.db.WithContext(ctx).Raw(`
		SELECT DISTINCT u.* FROM users u
		INNER JOIN friend_requests fr ON
			(fr.sender_id = u.id AND fr.receiver_id = ? AND fr.status = ?)
			OR (fr.receiver_id = u.id AND fr.sender_id = ? AND fr.status = ?)
		WHERE u.deleted_at IS NULL
	`, userID, models.StatusAccepted, userID, models.StatusAccepted).Scan(&friends).Error

	if err != nil {
		return nil, err
	}
	return friends, nil
}

func (r *friendRequestRepository) FindByUsers(ctx context.Context, userID1, userID2 uint) (*models.FriendRequest, error) {
	var friendRequest models.FriendRequest
	err := r.db.WithContext(ctx).Preload(preloadSender).Preload(preloadReceiver).
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			userID1, userID2, userID2, userID1).
		First(&friendRequest).Error
	if err != nil {
		return nil, err
	}
	return &friendRequest, nil
}
