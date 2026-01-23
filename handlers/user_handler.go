package handlers

import (
	"crypto/sha256"
	"fmt"
	"strconv"

	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/mappers"
	"github.com/PI-Team04-GameClub/gameclub-backend/repositories"
	"github.com/PI-Team04-GameClub/gameclub-backend/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserHandler struct {
	userRepo repositories.UserRepository
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{userRepo: repositories.NewUserRepository(db)}
}

func NewUserHandlerWithRepo(userRepo repositories.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

func hashUserPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}

func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.userRepo.FindAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.InternalServerError("Failed to fetch users"))
	}

	return c.JSON(mappers.ToUserResponseList(users))
}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.BadRequest("Invalid user ID"))
	}

	user, err := h.userRepo.FindByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.NotFound())
	}

	return c.JSON(mappers.ToUserResponse(user))
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req dtos.CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.BadRequest("Invalid request body"))
	}

	if req.Email == "" || req.Password == "" || req.FirstName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.BadRequest("Email, password, and first name are required"))
	}

	if len(req.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(utils.BadRequest("Password must be at least 6 characters"))
	}

	existingUser, _ := h.userRepo.FindByEmail(c.Context(), req.Email)
	if existingUser != nil {
		return c.Status(fiber.StatusConflict).JSON(utils.Conflict("Email already registered"))
	}

	user := mappers.ToUserModel(req)
	user.Password = hashUserPassword(req.Password)

	if err := h.userRepo.Create(c.Context(), &user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.InternalServerError("Failed to create user"))
	}

	return c.Status(fiber.StatusCreated).JSON(mappers.ToUserResponse(&user))
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.BadRequest("Invalid user ID"))
	}

	user, err := h.userRepo.FindByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.NotFound())
	}

	var req dtos.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.BadRequest("Invalid request body"))
	}

	if req.Email != "" && req.Email != user.Email {
		existingUser, _ := h.userRepo.FindByEmail(c.Context(), req.Email)
		if existingUser != nil {
			return c.Status(fiber.StatusConflict).JSON(utils.Conflict("Email already registered"))
		}
	}

	if req.Password != "" && len(req.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(utils.BadRequest("Password must be at least 6 characters"))
	}

	updatedUser := mappers.UpdateUserFromRequest(user, req)
	if req.Password != "" {
		updatedUser.Password = hashUserPassword(req.Password)
	}

	if err := h.userRepo.Update(c.Context(), updatedUser); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.InternalServerError("Failed to update user"))
	}

	return c.JSON(mappers.ToUserResponse(updatedUser))
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.BadRequest("Invalid user ID"))
	}

	user, err := h.userRepo.FindByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.NotFound())
	}

	if err := h.userRepo.Delete(c.Context(), user.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.InternalServerError("Failed to delete user"))
	}

	return c.SendStatus(fiber.StatusNoContent)
}
