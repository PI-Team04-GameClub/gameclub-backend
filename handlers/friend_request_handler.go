package handlers

import (
	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/mappers"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/PI-Team04-GameClub/gameclub-backend/repositories"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type FriendRequestHandler struct {
	friendRequestRepo repositories.FriendRequestRepository
	userRepo          repositories.UserRepository
}

func NewFriendRequestHandler(db *gorm.DB) *FriendRequestHandler {
	return &FriendRequestHandler{
		friendRequestRepo: repositories.NewFriendRequestRepository(db),
		userRepo:          repositories.NewUserRepository(db),
	}
}

func NewFriendRequestHandlerWithRepo(friendRequestRepo repositories.FriendRequestRepository, userRepo repositories.UserRepository) *FriendRequestHandler {
	return &FriendRequestHandler{
		friendRequestRepo: friendRequestRepo,
		userRepo:          userRepo,
	}
}

func (h *FriendRequestHandler) CreateFriendRequest(c *fiber.Ctx) error {
	var req dtos.CreateFriendRequestRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.SenderID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Sender ID is required"})
	}

	if req.ReceiverID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Receiver ID is required"})
	}

	// Cannot send friend request to self
	if req.SenderID == req.ReceiverID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot send friend request to yourself"})
	}

	// Verify sender exists
	_, err := h.userRepo.FindByID(c.Context(), req.SenderID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Sender not found"})
	}

	// Verify receiver exists
	_, err = h.userRepo.FindByID(c.Context(), req.ReceiverID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Receiver not found"})
	}

	// Check if friend request already exists between these users
	existingRequest, _ := h.friendRequestRepo.FindByUsers(c.Context(), req.SenderID, req.ReceiverID)
	if existingRequest != nil {
		if existingRequest.Status == models.StatusPending {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Friend request already pending"})
		}
		if existingRequest.Status == models.StatusAccepted {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Already friends"})
		}
	}

	friendRequest := mappers.ToFriendRequestModel(req)
	if err := h.friendRequestRepo.Create(c.Context(), &friendRequest); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create friend request"})
	}

	// Fetch the created friend request with relations
	createdFriendRequest, err := h.friendRequestRepo.FindByID(c.Context(), friendRequest.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve created friend request"})
	}

	return c.Status(fiber.StatusCreated).JSON(mappers.ToFriendRequestResponse(createdFriendRequest))
}

func (h *FriendRequestHandler) AcceptFriendRequest(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid friend request ID"})
	}

	friendRequest, err := h.friendRequestRepo.FindByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Friend request not found"})
	}

	if friendRequest.Status != models.StatusPending {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Friend request is not pending"})
	}

	friendRequest.Status = models.StatusAccepted
	if err := h.friendRequestRepo.Update(c.Context(), friendRequest); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to accept friend request"})
	}

	updatedFriendRequest, err := h.friendRequestRepo.FindByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve updated friend request"})
	}

	return c.JSON(mappers.ToFriendRequestResponse(updatedFriendRequest))
}

func (h *FriendRequestHandler) DeclineFriendRequest(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid friend request ID"})
	}

	friendRequest, err := h.friendRequestRepo.FindByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Friend request not found"})
	}

	if friendRequest.Status != models.StatusPending {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Friend request is not pending"})
	}

	friendRequest.Status = models.StatusDeclined
	if err := h.friendRequestRepo.Update(c.Context(), friendRequest); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decline friend request"})
	}

	updatedFriendRequest, err := h.friendRequestRepo.FindByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve updated friend request"})
	}

	return c.JSON(mappers.ToFriendRequestResponse(updatedFriendRequest))
}

func (h *FriendRequestHandler) DeleteFriendRequest(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid friend request ID"})
	}

	_, err = h.friendRequestRepo.FindByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Friend request not found"})
	}

	if err := h.friendRequestRepo.Delete(c.Context(), uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete friend request"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *FriendRequestHandler) GetSentFriendRequests(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Verify user exists
	_, err = h.userRepo.FindByID(c.Context(), uint(userID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	friendRequests, err := h.friendRequestRepo.FindBySenderID(c.Context(), uint(userID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve friend requests"})
	}

	return c.JSON(mappers.ToFriendRequestResponseList(friendRequests))
}

func (h *FriendRequestHandler) GetReceivedFriendRequests(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Verify user exists
	_, err = h.userRepo.FindByID(c.Context(), uint(userID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	friendRequests, err := h.friendRequestRepo.FindPendingByReceiverID(c.Context(), uint(userID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve friend requests"})
	}

	return c.JSON(mappers.ToFriendRequestResponseList(friendRequests))
}

func (h *FriendRequestHandler) GetFriends(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Verify user exists
	_, err = h.userRepo.FindByID(c.Context(), uint(userID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	friends, err := h.friendRequestRepo.FindFriendsByUserID(c.Context(), uint(userID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve friends"})
	}

	return c.JSON(mappers.ToFriendResponseList(friends))
}
