package handlers

import (
	"github.com/PI-Team04-GameClub/gameclub-backend/db"
	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/mappers"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/gofiber/fiber/v2"
)

func GetAllGames(c *fiber.Ctx) error {
	var games []models.Game

	if err := db.DB.Find(&games).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch games",
		})
	}

	return c.JSON(mappers.ToGameResponseList(games))
}

func GetGameByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var game models.Game
	if err := db.DB.First(&game, id).Error; err != nil {
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

	// Use Builder pattern through mapper
	game := mappers.ToGameModel(req)

	if err := db.DB.Create(&game).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create game",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(mappers.ToGameResponse(&game))
}

func UpdateGame(c *fiber.Ctx) error {
	id := c.Params("id")

	var game models.Game
	if err := db.DB.First(&game, id).Error; err != nil {
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

	// Use Builder pattern through mapper to update the game
	updatedGame := mappers.UpdateGameFromRequest(&game, req)

	if err := db.DB.Save(updatedGame).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update game",
		})
	}

	return c.JSON(mappers.ToGameResponse(updatedGame))
}

func DeleteGame(c *fiber.Ctx) error {
	id := c.Params("id")

	var game models.Game
	if err := db.DB.First(&game, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Game not found",
		})
	}

	if err := db.DB.Delete(&game).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete game",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
