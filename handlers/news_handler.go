package handlers

import (
	"strconv"

	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/mappers"
	"github.com/PI-Team04-GameClub/gameclub-backend/repositories"
	"github.com/PI-Team04-GameClub/gameclub-backend/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type NewsHandler struct {
	newsRepo repositories.NewsRepository
	userRepo repositories.UserRepository
}

func NewNewsHandler(db *gorm.DB) *NewsHandler {
	return &NewsHandler{
		newsRepo: repositories.NewNewsRepository(db),
		userRepo: repositories.NewUserRepository(db),
	}
}

func NewNewsHandlerWithRepo(newsRepo repositories.NewsRepository, userRepo repositories.UserRepository) *NewsHandler {
	return &NewsHandler{
		newsRepo: newsRepo,
		userRepo: userRepo,
	}
}

func (h *NewsHandler) GetNews(c *fiber.Ctx) error {
	newsList, err := h.newsRepo.FindAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.InternalServerError("Failed to fetch news"))
	}

	responses := mappers.ToNewsResponseList(newsList)
	return c.JSON(responses)
}

func (h *NewsHandler) CreateNews(c *fiber.Ctx) error {
	ctx := c.Context()

	var req dtos.CreateNewsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.BadRequest("Invalid request body"))
	}

	if req.AuthorId == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(utils.BadRequest("Author ID is required"))
	}

	_, err := h.userRepo.FindByID(ctx, req.AuthorId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.BadRequest("Invalid author ID"))
	}

	news := mappers.ToNewsModel(req, req.AuthorId)
	if err := h.newsRepo.Create(ctx, &news); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.InternalServerError("Failed to create news"))
	}

	createdNews, _ := h.newsRepo.FindByID(ctx, int(news.ID))

	response := mappers.ToNewsResponse(createdNews, createdNews.Author.FirstName)
	return c.Status(fiber.StatusCreated).JSON(response)
}

func (h *NewsHandler) UpdateNews(c *fiber.Ctx) error {
	ctx := c.Context()

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.BadRequest("Invalid news ID"))
	}

	news, err := h.newsRepo.FindByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.NotFound())
	}

	var req dtos.CreateNewsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.BadRequest("Invalid request body"))
	}

	news.Title = req.Title
	news.Description = req.Description

	if err := h.newsRepo.Update(ctx, news); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.InternalServerError("Failed to update news"))
	}

	response := mappers.ToNewsResponse(news, news.Author.FirstName)
	return c.JSON(response)
}

func (h *NewsHandler) DeleteNews(c *fiber.Ctx) error {
	ctx := c.Context()

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.BadRequest("Invalid news ID"))
	}

	_, err = h.newsRepo.FindByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.NotFound())
	}

	if err := h.newsRepo.Delete(ctx, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.InternalServerError("Failed to delete news"))
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
