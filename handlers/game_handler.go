package handlers

import (
	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/mappers"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type GameHandler struct {
	db *gorm.DB
}

func NewGameHandler(db *gorm.DB) *GameHandler {
	return &GameHandler{db: db}
}

func (h *GameHandler) GetAllGames(c *fiber.Ctx) error {
	games, err := gorm.G[models.Game](h.db).Find(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch games",
		})
	}

	return c.JSON(mappers.ToGameResponseList(games))
}

func (h *GameHandler) GetGameByID(c *fiber.Ctx) error {
	id := c.Params("id")

	game, err := gorm.G[models.Game](h.db).Where("id = ?", id).First(c.Context())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Game not found",
		})
	}

	return c.JSON(mappers.ToGameResponse(&game))
}

func (h *GameHandler) CreateGame(c *fiber.Ctx) error {
	var req dtos.CreateGameRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	game := mappers.ToGameModel(req)

	if err := gorm.G[models.Game](h.db).Create(c.Context(), &game); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create game",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(mappers.ToGameResponse(&game))
}

func (h *GameHandler) UpdateGame(c *fiber.Ctx) error {
	id := c.Params("id")

	game, err := gorm.G[models.Game](h.db).Where("id = ?", id).First(c.Context())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Game not found",
		})
	}

	var req dtos.CreateGameRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	updatedGame := mappers.UpdateGameFromRequest(&game, req)

	if _, err := gorm.G[models.Game](h.db).Where("id = ?", game.ID).Updates(c.Context(), *updatedGame); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update game",
		})
	}

	return c.JSON(mappers.ToGameResponse(updatedGame))
}

func (h *GameHandler) DeleteGame(c *fiber.Ctx) error {
	id := c.Params("id")

	game, err := gorm.G[models.Game](h.db).Where("id = ?", id).First(c.Context())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Game not found",
		})
	}

	if _, err := gorm.G[models.Game](h.db).Where("id = ?", game.ID).Delete(c.Context()); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete game",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
