package handlers

import (
	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/mappers"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TeamHandler struct {
	db *gorm.DB
}

func NewTeamHandler(db *gorm.DB) *TeamHandler {
	return &TeamHandler{db: db}
}

func (h *TeamHandler) GetAllTeams(c *fiber.Ctx) error {
	var teams []models.Team

	if err := h.db.Find(&teams).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch teams",
		})
	}

	return c.JSON(mappers.ToTeamResponseList(teams))
}

func (h *TeamHandler) GetTeamByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var team models.Team
	if err := h.db.First(&team, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Team not found",
		})
	}

	return c.JSON(mappers.ToTeamResponse(&team))
}

func (h *TeamHandler) CreateTeam(c *fiber.Ctx) error {
	var req dtos.CreateTeamRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	team := mappers.ToTeamModel(req)

	if err := h.db.Create(&team).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create team",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(mappers.ToTeamResponse(&team))
}

func (h *TeamHandler) UpdateTeam(c *fiber.Ctx) error {
	id := c.Params("id")

	var team models.Team
	if err := h.db.First(&team, id).Error; err != nil {
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

	if err := h.db.Save(&team).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update team",
		})
	}

	return c.JSON(mappers.ToTeamResponse(&team))
}

func (h *TeamHandler) DeleteTeam(c *fiber.Ctx) error {
	id := c.Params("id")

	var team models.Team
	if err := h.db.First(&team, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Team not found",
		})
	}

	if err := h.db.Delete(&team).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete team",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
