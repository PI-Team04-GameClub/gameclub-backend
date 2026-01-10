package handlers

import (
	"context"
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
	tournaments, err := gorm.G[models.Tournament](h.db).Preload("Game", nil).Find(c.Context())
	if err != nil {
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

	tournament, err := gorm.G[models.Tournament](h.db).Preload("Game", nil).Where("id = ?", id).First(c.Context())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Tournament not found",
		})
	}

	response := mappers.ToTournamentResponse(&tournament)
	return c.JSON(response)
}

func (h *TournamentHandler) CreateTournament(c *fiber.Ctx) error {
	ctx := c.Context()

	var req dtos.CreateTournamentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	_, err := gorm.G[models.Game](h.db).Where("id = ?", req.GameId).First(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Game not found",
		})
	}

	tournament := mappers.ToTournamentModel(req)

	if err := gorm.G[models.Tournament](h.db).Create(ctx, &tournament); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create tournament",
		})
	}

	createdTournament, _ := gorm.G[models.Tournament](h.db).Preload("Game", nil).Where("id = ?", tournament.ID).First(ctx)

	users, _ := gorm.G[models.User](h.db).Find(context.Background())

	userEmails := make(map[string]string)
	for _, user := range users {
		userEmails[user.Email] = user.FirstName
	}

	createdTournament.Attach(observer.NewEmailNotifier(userEmails))
	createdTournament.Attach(observer.NewLogNotifier())

	createdTournament.NotifyCreated()

	response := mappers.ToTournamentResponse(&createdTournament)
	return c.Status(fiber.StatusCreated).JSON(response)
}

func (h *TournamentHandler) UpdateTournament(c *fiber.Ctx) error {
	ctx := c.Context()

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tournament ID",
		})
	}

	tournament, err := gorm.G[models.Tournament](h.db).Where("id = ?", id).First(ctx)
	if err != nil {
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

	_, err = gorm.G[models.Game](h.db).Where("id = ?", req.GameId).First(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Game not found",
		})
	}

	updatedTournament := mappers.UpdateTournamentFromRequest(&tournament, req)

	if _, err := gorm.G[models.Tournament](h.db).Where("id = ?", tournament.ID).Updates(ctx, *updatedTournament); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update tournament",
		})
	}

	result, _ := gorm.G[models.Tournament](h.db).Preload("Game", nil).Where("id = ?", updatedTournament.ID).First(ctx)

	response := mappers.ToTournamentResponse(&result)
	return c.JSON(response)
}

func (h *TournamentHandler) DeleteTournament(c *fiber.Ctx) error {
	ctx := c.Context()

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tournament ID",
		})
	}

	tournament, err := gorm.G[models.Tournament](h.db).Where("id = ?", id).First(ctx)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Tournament not found",
		})
	}

	if _, err := gorm.G[models.Tournament](h.db).Where("id = ?", tournament.ID).Delete(ctx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete tournament",
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
