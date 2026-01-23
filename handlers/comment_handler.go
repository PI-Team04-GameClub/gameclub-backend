package handlers

import (
	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/mappers"
	"github.com/PI-Team04-GameClub/gameclub-backend/repositories"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CommentHandler struct {
	commentRepo repositories.CommentRepository
	userRepo    repositories.UserRepository
	newsRepo    repositories.NewsRepository
}

func NewCommentHandler(db *gorm.DB) *CommentHandler {
	return &CommentHandler{
		commentRepo: repositories.NewCommentRepository(db),
		userRepo:    repositories.NewUserRepository(db),
		newsRepo:    repositories.NewNewsRepository(db),
	}
}

func NewCommentHandlerWithRepo(commentRepo repositories.CommentRepository, userRepo repositories.UserRepository, newsRepo repositories.NewsRepository) *CommentHandler {
	return &CommentHandler{
		commentRepo: commentRepo,
		userRepo:    userRepo,
		newsRepo:    newsRepo,
	}
}

func (h *CommentHandler) CreateComment(c *fiber.Ctx) error {
	var req dtos.CreateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Content is required"})
	}

	if req.UserID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID is required"})
	}

	if req.NewsID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "News ID is required"})
	}

	// Verify user exists
	_, err := h.userRepo.FindByID(c.Context(), req.UserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Verify news exists
	_, err = h.newsRepo.FindByID(c.Context(), int(req.NewsID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "News not found"})
	}

	comment := mappers.ToCommentModel(req)
	if err := h.commentRepo.Create(c.Context(), &comment); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create comment"})
	}

	// Fetch the created comment with relations
	createdComment, err := h.commentRepo.FindByID(c.Context(), comment.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve created comment"})
	}

	return c.Status(fiber.StatusCreated).JSON(mappers.ToCommentResponse(createdComment))
}

func (h *CommentHandler) UpdateComment(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid comment ID"})
	}

	var req dtos.UpdateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Content is required"})
	}

	comment, err := h.commentRepo.FindByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Comment not found"})
	}

	comment.Content = req.Content

	if err := h.commentRepo.Update(c.Context(), comment); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update comment"})
	}

	// Fetch updated comment with relations
	updatedComment, err := h.commentRepo.FindByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve updated comment"})
	}

	return c.JSON(mappers.ToCommentResponse(updatedComment))
}

func (h *CommentHandler) GetCommentsByUserID(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Verify user exists
	_, err = h.userRepo.FindByID(c.Context(), uint(userID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	comments, err := h.commentRepo.FindByUserID(c.Context(), uint(userID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve comments"})
	}

	return c.JSON(mappers.ToCommentResponseList(comments))
}

func (h *CommentHandler) GetCommentsByNewsID(c *fiber.Ctx) error {
	newsID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid news ID"})
	}

	// Verify news exists
	_, err = h.newsRepo.FindByID(c.Context(), newsID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "News not found"})
	}

	comments, err := h.commentRepo.FindByNewsID(c.Context(), uint(newsID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve comments"})
	}

	return c.JSON(mappers.ToCommentResponseList(comments))
}
