package handlers

import (
	"strconv"

	"github.com/PI-Team04-GameClub/gameclub-backend/db"
	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/mappers"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/gofiber/fiber/v2"
)

func GetTournaments(c *fiber.Ctx) error {
	var tournaments []models.Tournament
	if err := db.DB.Preload("Game").Find(&tournaments).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch tournaments",
		})
	}

	responses := mappers.ToTournamentResponseList(tournaments)
	return c.JSON(responses)
}

func GetTournamentByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tournament ID",
		})
	}

	var tournament models.Tournament
	if err := db.DB.Preload("Game").First(&tournament, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Tournament not found",
		})
	}

	response := mappers.ToTournamentResponse(&tournament)
	return c.JSON(response)
}

func CreateTournament(c *fiber.Ctx) error {
	var req dtos.CreateTournamentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	tournament := mappers.ToTournamentModel(req)
	if err := db.DB.Create(&tournament).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create tournament",
		})
	}

	// Preload Game to get the game name for response
	db.DB.Preload("Game").First(&tournament, tournament.ID)

	response := mappers.ToTournamentResponse(&tournament)
	return c.Status(fiber.StatusCreated).JSON(response)
}

func UpdateTournament(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tournament ID",
		})
	}

	var tournament models.Tournament
	if err := db.DB.First(&tournament, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Tournament not found",
		})
	}

	var req dtos.CreateTournamentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	tournament.Name = req.Name
	tournament.GameID = req.GameId
	tournament.PrizePool = req.PrizePool
	tournament.StartDate = req.StartDate

	if err := db.DB.Save(&tournament).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update tournament",
		})
	}

	// Preload Game to get the game name for response
	db.DB.Preload("Game").First(&tournament, tournament.ID)

	response := mappers.ToTournamentResponse(&tournament)
	return c.JSON(response)
}

func DeleteTournament(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tournament ID",
		})
	}

	var tournament models.Tournament
	if err := db.DB.First(&tournament, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Tournament not found",
		})
	}

	if err := db.DB.Delete(&tournament).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete tournament",
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
