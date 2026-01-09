package handlers

import (
	"github.com/PI-Team04-GameClub/gameclub-backend/db"
	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/mappers"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetAllTeams(c *fiber.Ctx) error {
	teams, err := gorm.G[models.Team](db.DB).Find(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch teams",
		})
	}

	return c.JSON(mappers.ToTeamResponseList(teams))
}

func GetTeamByID(c *fiber.Ctx) error {
	id := c.Params("id")

	team, err := gorm.G[models.Team](db.DB).Where("id = ?", id).First(c.Context())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Team not found",
		})
	}

	return c.JSON(mappers.ToTeamResponse(&team))
}

func CreateTeam(c *fiber.Ctx) error {
	var req dtos.CreateTeamRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	team := mappers.ToTeamModel(req)

	if err := gorm.G[models.Team](db.DB).Create(c.Context(), &team); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create team",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(mappers.ToTeamResponse(&team))
}

func UpdateTeam(c *fiber.Ctx) error {
	id := c.Params("id")

	team, err := gorm.G[models.Team](db.DB).Where("id = ?", id).First(c.Context())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Team not found",
		})
	}

	var req dtos.CreateTeamRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	team.Name = req.Name

	if _, err := gorm.G[models.Team](db.DB).Where("id = ?", team.ID).Updates(c.Context(), team); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update team",
		})
	}

	return c.JSON(mappers.ToTeamResponse(&team))
}

func DeleteTeam(c *fiber.Ctx) error {
	id := c.Params("id")

	team, err := gorm.G[models.Team](db.DB).Where("id = ?", id).First(c.Context())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Team not found",
		})
	}

	if _, err := gorm.G[models.Team](db.DB).Where("id = ?", team.ID).Delete(c.Context()); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete team",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
