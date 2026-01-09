package handlers

import (
	"github.com/PI-Team04-GameClub/gameclub-backend/db"
	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/mappers"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetAllGames(c *fiber.Ctx) error {
	games, err := gorm.G[models.Game](db.DB).Find(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch games",
		})
	}

	return c.JSON(mappers.ToGameResponseList(games))
}

func GetGameByID(c *fiber.Ctx) error {
	id := c.Params("id")

	game, err := gorm.G[models.Game](db.DB).Where("id = ?", id).First(c.Context())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Game not found",
		})
	}

	return c.JSON(mappers.ToGameResponse(&game))
}

func CreateGame(c *fiber.Ctx) error {
	var req dtos.CreateGameRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	game := mappers.ToGameModel(req)

	if err := gorm.G[models.Game](db.DB).Create(c.Context(), &game); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create game",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(mappers.ToGameResponse(&game))
}

func UpdateGame(c *fiber.Ctx) error {
	id := c.Params("id")

	game, err := gorm.G[models.Game](db.DB).Where("id = ?", id).First(c.Context())
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

	if _, err := gorm.G[models.Game](db.DB).Where("id = ?", game.ID).Updates(c.Context(), *updatedGame); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update game",
		})
	}

	return c.JSON(mappers.ToGameResponse(updatedGame))
}

func DeleteGame(c *fiber.Ctx) error {
	id := c.Params("id")

	game, err := gorm.G[models.Game](db.DB).Where("id = ?", id).First(c.Context())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Game not found",
		})
	}

	if _, err := gorm.G[models.Game](db.DB).Where("id = ?", game.ID).Delete(c.Context()); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete game",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
