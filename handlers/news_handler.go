package handlers

import (
	"strconv"

	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/mappers"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type NewsHandler struct {
	db *gorm.DB
}

func NewNewsHandler(db *gorm.DB) *NewsHandler {
	return &NewsHandler{db: db}
}

func (h *NewsHandler) GetNews(c *fiber.Ctx) error {
	var newsList []models.News
	if err := h.db.Preload("Author").Find(&newsList).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch news",
		})
	}

	responses := mappers.ToNewsResponseList(newsList)
	return c.JSON(responses)
}

func (h *NewsHandler) CreateNews(c *fiber.Ctx) error {
	var req dtos.CreateNewsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.AuthorId == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Author ID is required",
		})
	}

	var author models.User
	if err := h.db.First(&author, req.AuthorId).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid author ID",
		})
	}

	news := mappers.ToNewsModel(req, req.AuthorId)
	if err := h.db.Create(&news).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create news",
		})
	}

	h.db.Preload("Author").First(&news, news.ID)

	response := mappers.ToNewsResponse(&news, news.Author.FirstName)
	return c.Status(fiber.StatusCreated).JSON(response)
}

func (h *NewsHandler) UpdateNews(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid news ID",
		})
	}

	var news models.News
	if err := h.db.Preload("Author").First(&news, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "News not found",
		})
	}

	var req dtos.CreateNewsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	news.Title = req.Title
	news.Description = req.Description

	if err := h.db.Save(&news).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update news",
		})
	}

	response := mappers.ToNewsResponse(&news, news.Author.FirstName)
	return c.JSON(response)
}

func (h *NewsHandler) DeleteNews(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid news ID",
		})
	}

	var news models.News
	if err := h.db.First(&news, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "News not found",
		})
	}

	if err := h.db.Delete(&news).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete news",
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
