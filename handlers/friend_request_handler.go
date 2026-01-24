package handlers

import (
	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/mappers"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/PI-Team04-GameClub/gameclub-backend/repositories"
	"github.com/PI-Team04-GameClub/gameclub-backend/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	errInvalidFriendRequestID         = "Invalid friend request ID"
	errFriendRequestNotFound          = "Friend request not found"
	errFriendRequestNotPending        = "Friend request is not pending"
	errInvalidUserID                  = "Invalid user ID"
	errUserNotFound                   = "User not found"
	errFailedToRetrieveFriendRequests = "Failed to retrieve friend requests"
	errFailedToRetrieveUpdated        = "Failed to retrieve updated friend request"
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
		return c.Status(fiber.StatusBadRequest).JSON(utils.BadRequest("Invalid request body"))
	}

	if req.SenderID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(utils.BadRequest("Sender ID is required"))
	}

	if req.ReceiverID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(utils.BadRequest("Receiver ID is required"))
	}

	// Cannot send friend request to self
	if req.SenderID == req.ReceiverID {
		return c.Status(fiber.StatusBadRequest).JSON(utils.BadRequest("Cannot send friend request to yourself"))
	}

	// Verify sender exists
	_, err := h.userRepo.FindByID(c.Context(), req.SenderID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.BadRequest("Sender not found"))
	}

	// Verify receiver exists
	_, err = h.userRepo.FindByID(c.Context(), req.ReceiverID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.BadRequest("Receiver not found"))
	}

	// Check if friend request already exists between these users
	existingRequest, _ := h.friendRequestRepo.FindByUsers(c.Context(), req.SenderID, req.ReceiverID)
	if existingRequest != nil {
		if existingRequest.Status == models.StatusPending {
			return c.Status(fiber.StatusConflict).JSON(utils.Conflict("Friend request already pending"))
		}
		if existingRequest.Status == models.StatusAccepted {
			return c.Status(fiber.StatusConflict).JSON(utils.Conflict("Already friends"))
		}
	}

	friendRequest := mappers.ToFriendRequestModel(req)
	if err := h.friendRequestRepo.Create(c.Context(), &friendRequest); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.InternalServerError("Failed to create friend request"))
	}

	// Fetch the created friend request with relations
	createdFriendRequest, err := h.friendRequestRepo.FindByID(c.Context(), friendRequest.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.InternalServerError("Failed to retrieve created friend request"))
	}

	return c.Status(fiber.StatusCreated).JSON(mappers.ToFriendRequestResponse(createdFriendRequest))
}

func (h *FriendRequestHandler) AcceptFriendRequest(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.BadRequest(errInvalidFriendRequestID))
	}

	friendRequest, err := h.friendRequestRepo.FindByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.BadRequest(errFriendRequestNotFound))
	}

	if friendRequest.Status != models.StatusPending {
		return c.Status(fiber.StatusBadRequest).JSON(utils.BadRequest(errFriendRequestNotPending))
	}

	friendRequest.Status = models.StatusAccepted
	if err := h.friendRequestRepo.Update(c.Context(), friendRequest); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.InternalServerError("Failed to accept friend request"))
	}

	updatedFriendRequest, err := h.friendRequestRepo.FindByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.InternalServerError(errFailedToRetrieveUpdated))
	}

	return c.JSON(mappers.ToFriendRequestResponse(updatedFriendRequest))
}

func (h *FriendRequestHandler) DeclineFriendRequest(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.BadRequest(errInvalidFriendRequestID))
	}

	friendRequest, err := h.friendRequestRepo.FindByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.BadRequest(errFriendRequestNotFound))
	}

	if friendRequest.Status != models.StatusPending {
		return c.Status(fiber.StatusBadRequest).JSON(utils.BadRequest(errFriendRequestNotPending))
	}

	friendRequest.Status = models.StatusDeclined
	if err := h.friendRequestRepo.Update(c.Context(), friendRequest); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.InternalServerError("Failed to decline friend request"))
	}

	updatedFriendRequest, err := h.friendRequestRepo.FindByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.InternalServerError(errFailedToRetrieveUpdated))
	}

	return c.JSON(mappers.ToFriendRequestResponse(updatedFriendRequest))
}

func (h *FriendRequestHandler) DeleteFriendRequest(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.BadRequest(errInvalidFriendRequestID))
	}

	_, err = h.friendRequestRepo.FindByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.BadRequest(errFriendRequestNotFound))
	}

	if err := h.friendRequestRepo.Delete(c.Context(), uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.InternalServerError("Failed to delete friend request"))
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *FriendRequestHandler) GetSentFriendRequests(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.BadRequest(errInvalidUserID))
	}

	// Verify user exists
	_, err = h.userRepo.FindByID(c.Context(), uint(userID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.BadRequest(errUserNotFound))
	}

	friendRequests, err := h.friendRequestRepo.FindBySenderID(c.Context(), uint(userID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.InternalServerError(errFailedToRetrieveFriendRequests))
	}

	return c.JSON(mappers.ToFriendRequestResponseList(friendRequests))
}

func (h *FriendRequestHandler) GetReceivedFriendRequests(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.BadRequest(errInvalidUserID))
	}

	// Verify user exists
	_, err = h.userRepo.FindByID(c.Context(), uint(userID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.BadRequest(errUserNotFound))
	}

	friendRequests, err := h.friendRequestRepo.FindPendingByReceiverID(c.Context(), uint(userID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.InternalServerError(errFailedToRetrieveFriendRequests))
	}

	return c.JSON(mappers.ToFriendRequestResponseList(friendRequests))
}

func (h *FriendRequestHandler) GetFriends(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.BadRequest(errInvalidUserID))
	}

	// Verify user exists
	_, err = h.userRepo.FindByID(c.Context(), uint(userID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.BadRequest(errUserNotFound))
	}

	friends, err := h.friendRequestRepo.FindFriendsByUserID(c.Context(), uint(userID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.InternalServerError("Failed to retrieve friends"))
	}

	return c.JSON(mappers.ToFriendResponseList(friends))
}
