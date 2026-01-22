package handlers

import (
	"crypto/sha256"
	"fmt"

	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/PI-Team04-GameClub/gameclub-backend/repositories"
	"github.com/PI-Team04-GameClub/gameclub-backend/security"
	"github.com/PI-Team04-GameClub/gameclub-backend/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// errResponseSent is a sentinel error indicating that a response has already been sent
var errResponseSent = fmt.Errorf("response already sent")

type AuthHandler struct {
	userRepo repositories.UserRepository
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{userRepo: repositories.NewUserRepository(db)}
}

func NewAuthHandlerWithRepo(userRepo repositories.UserRepository) *AuthHandler {
	return &AuthHandler{userRepo: userRepo}
}

// Password utilities

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}

func verifyPassword(hashedPassword, password string) bool {
	return hashPassword(password) == hashedPassword
}

// Validation helpers

func validateRegisterRequest(req *dtos.RegisterRequest) error {
	if req.Email == "" || req.Password == "" || req.FirstName == "" {
		return fmt.Errorf("email, password, and first name are required")
	}
	if len(req.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}
	return nil
}

func validateLoginRequest(req *dtos.LoginRequest) error {
	if req.Email == "" || req.Password == "" {
		return fmt.Errorf("email and password are required")
	}
	return nil
}

// Response builders

func buildAuthResponse(user *models.User, token string) dtos.AuthResponse {
	return dtos.AuthResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Token:     token,
	}
}

func buildUserResponse(user *models.User) fiber.Map {
	return fiber.Map{
		"id":         user.ID,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
	}
}

// Handler methods

func (ah *AuthHandler) Register(c *fiber.Ctx) error {
	req, err := ah.parseRegisterRequest(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.BadRequest("Invalid request body"))
	}

	if err := validateRegisterRequest(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.BadRequest(err.Error()))
	}

	if err := ah.checkEmailAvailable(c); err != nil {
		return nil
	}

	user, err := ah.createUser(c, req)
	if err != nil {
		return nil
	}

	return ah.respondWithAuthToken(c, user, fiber.StatusCreated)
}

func (ah *AuthHandler) Login(c *fiber.Ctx) error {
	req, err := ah.parseLoginRequest(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.BadRequest("Invalid request body"))
	}

	if err := validateLoginRequest(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.BadRequest(err.Error()))
	}

	user, err := ah.findUserByEmail(c, req.Email)
	if err != nil {
		return nil
	}

	if !verifyPassword(user.Password, req.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.Unauthorized("Invalid password"))
	}

	return ah.respondWithAuthToken(c, user, fiber.StatusOK)
}

func (ah *AuthHandler) GetCurrentUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	return c.Status(fiber.StatusOK).JSON(buildUserResponse(user))
}

// Private helper methods

func (ah *AuthHandler) parseRegisterRequest(c *fiber.Ctx) (*dtos.RegisterRequest, error) {
	var req dtos.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (ah *AuthHandler) parseLoginRequest(c *fiber.Ctx) (*dtos.LoginRequest, error) {
	var req dtos.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (ah *AuthHandler) checkEmailAvailable(c *fiber.Ctx) error {
	var req dtos.RegisterRequest
	c.BodyParser(&req)
	ctx := c.Context()
	_, err := ah.userRepo.FindByEmail(ctx, req.Email)
	if err == nil {
		c.Status(fiber.StatusConflict).JSON(utils.Conflict("Email already registered"))
		return errResponseSent
	}
	if err != gorm.ErrRecordNotFound {
		c.Status(fiber.StatusInternalServerError).JSON(utils.InternalServerError("Database error"))
		return errResponseSent
	}
	return nil
}

func (ah *AuthHandler) createUser(c *fiber.Ctx, req *dtos.RegisterRequest) (*models.User, error) {
	ctx := c.Context()
	user := &models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  hashPassword(req.Password),
	}

	if err := ah.userRepo.Create(ctx, user); err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(utils.InternalServerError("Failed to create user"))
		return nil, errResponseSent
	}
	return user, nil
}

func (ah *AuthHandler) findUserByEmail(c *fiber.Ctx, email string) (*models.User, error) {
	ctx := c.Context()
	user, err := ah.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Status(fiber.StatusUnauthorized).JSON(utils.Unauthorized("User with this email address not found"))
			return nil, errResponseSent
		}
		c.Status(fiber.StatusInternalServerError).JSON(utils.InternalServerError("Database access error"))
		return nil, errResponseSent
	}
	return user, nil
}

func (ah *AuthHandler) respondWithAuthToken(c *fiber.Ctx, user *models.User, status int) error {
	token, err := security.GenerateToken(user.ID, user.Email, user.FirstName, user.LastName)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(utils.InternalServerError("Failed to generate token"))
		return nil
	}
	return c.Status(status).JSON(buildAuthResponse(user, token))
}
