package handlers

import (
	"strconv"

	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/mappers"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/PI-Team04-GameClub/gameclub-backend/observer"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TournamentHandler struct {
	db *gorm.DB
}

func NewTournamentHandler(db *gorm.DB) *TournamentHandler {
	return &TournamentHandler{db: db}
}

func (h *TournamentHandler) GetTournaments(c *fiber.Ctx) error {
	var tournaments []models.Tournament
	if err := h.db.Preload("Game").Find(&tournaments).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch tournaments",
		})
	}

	responses := mappers.ToTournamentResponseList(tournaments)
	return c.JSON(responses)
}

func (h *TournamentHandler) GetTournamentByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tournament ID",
		})
	}

	var tournament models.Tournament
	if err := h.db.Preload("Game").First(&tournament, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Tournament not found",
		})
	}

	response := mappers.ToTournamentResponse(&tournament)
	return c.JSON(response)
}

func (h *TournamentHandler) CreateTournament(c *fiber.Ctx) error {
	var req dtos.CreateTournamentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var game models.Game
	if err := h.db.First(&game, req.GameId).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Game not found",
		})
	}

	tournament := mappers.ToTournamentModel(req)

	if err := h.db.Create(&tournament).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create tournament",
		})
	}

	h.db.Preload("Game").First(&tournament, tournament.ID)

	var users []models.User
	h.db.Find(&users)

	userEmails := make(map[string]string)
	for _, user := range users {
		userEmails[user.Email] = user.FirstName
	}

	tournament.Attach(observer.NewEmailNotifier(userEmails))
	tournament.Attach(observer.NewLogNotifier())

	tournament.NotifyCreated()

	response := mappers.ToTournamentResponse(&tournament)
	return c.Status(fiber.StatusCreated).JSON(response)
}

func (h *TournamentHandler) UpdateTournament(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tournament ID",
		})
	}

	var tournament models.Tournament
	if err := h.db.First(&tournament, id).Error; err != nil {
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

	var game models.Game
	if err := h.db.First(&game, req.GameId).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Game not found",
		})
	}

	updatedTournament := mappers.UpdateTournamentFromRequest(&tournament, req)

	if err := h.db.Save(updatedTournament).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update tournament",
		})
	}

	h.db.Preload("Game").First(updatedTournament, updatedTournament.ID)

	response := mappers.ToTournamentResponse(updatedTournament)
	return c.JSON(response)
}

func (h *TournamentHandler) DeleteTournament(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tournament ID",
		})
	}

	var tournament models.Tournament
	if err := h.db.First(&tournament, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Tournament not found",
		})
	}

	if err := h.db.Delete(&tournament).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete tournament",
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
