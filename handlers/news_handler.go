package handlers

import (
	"strconv"

	"github.com/PI-Team04-GameClub/gameclub-backend/db"
	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/mappers"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetNews(c *fiber.Ctx) error {
	newsList, err := gorm.G[models.News](db.DB).Preload("Author", nil).Find(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch news",
		})
	}

	responses := mappers.ToNewsResponseList(newsList)
	return c.JSON(responses)
}

func CreateNews(c *fiber.Ctx) error {
	ctx := c.Context()

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

	_, err := gorm.G[models.User](db.DB).Where("id = ?", req.AuthorId).First(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid author ID",
		})
	}

	news := mappers.ToNewsModel(req, req.AuthorId)
	if err := gorm.G[models.News](db.DB).Create(ctx, &news); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create news",
		})
	}

	createdNews, _ := gorm.G[models.News](db.DB).Preload("Author", nil).Where("id = ?", news.ID).First(ctx)

	response := mappers.ToNewsResponse(&createdNews, createdNews.Author.FirstName)
	return c.Status(fiber.StatusCreated).JSON(response)
}

func UpdateNews(c *fiber.Ctx) error {
	ctx := c.Context()

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid news ID",
		})
	}

	news, err := gorm.G[models.News](db.DB).Preload("Author", nil).Where("id = ?", id).First(ctx)
	if err != nil {
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

	if _, err := gorm.G[models.News](db.DB).Where("id = ?", news.ID).Updates(ctx, news); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update news",
		})
	}

	response := mappers.ToNewsResponse(&news, news.Author.FirstName)
	return c.JSON(response)
}

func DeleteNews(c *fiber.Ctx) error {
	ctx := c.Context()

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid news ID",
		})
	}

	news, err := gorm.G[models.News](db.DB).Where("id = ?", id).First(ctx)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "News not found",
		})
	}

	if _, err := gorm.G[models.News](db.DB).Where("id = ?", news.ID).Delete(ctx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete news",
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
