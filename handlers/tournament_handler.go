package handlers

import (
	"context"
	"strconv"

	"github.com/PI-Team04-GameClub/gameclub-backend/db"
	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/mappers"
	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/PI-Team04-GameClub/gameclub-backend/observer"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetTournaments(c *fiber.Ctx) error {
	tournaments, err := gorm.G[models.Tournament](db.DB).Preload("Game", nil).Find(c.Context())
	if err != nil {
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

	tournament, err := gorm.G[models.Tournament](db.DB).Preload("Game", nil).Where("id = ?", id).First(c.Context())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Tournament not found",
		})
	}

	response := mappers.ToTournamentResponse(&tournament)
	return c.JSON(response)
}

func CreateTournament(c *fiber.Ctx) error {
	ctx := c.Context()

	var req dtos.CreateTournamentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	_, err := gorm.G[models.Game](db.DB).Where("id = ?", req.GameId).First(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Game not found",
		})
	}

	tournament := mappers.ToTournamentModel(req)

	if err := gorm.G[models.Tournament](db.DB).Create(ctx, &tournament); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create tournament",
		})
	}

	createdTournament, _ := gorm.G[models.Tournament](db.DB).Preload("Game", nil).Where("id = ?", tournament.ID).First(ctx)

	users, _ := gorm.G[models.User](db.DB).Find(context.Background())

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

func UpdateTournament(c *fiber.Ctx) error {
	ctx := c.Context()

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tournament ID",
		})
	}

	tournament, err := gorm.G[models.Tournament](db.DB).Where("id = ?", id).First(ctx)
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

	_, err = gorm.G[models.Game](db.DB).Where("id = ?", req.GameId).First(ctx)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Game not found",
		})
	}

	updatedTournament := mappers.UpdateTournamentFromRequest(&tournament, req)

	if _, err := gorm.G[models.Tournament](db.DB).Where("id = ?", tournament.ID).Updates(ctx, *updatedTournament); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update tournament",
		})
	}

	result, _ := gorm.G[models.Tournament](db.DB).Preload("Game", nil).Where("id = ?", updatedTournament.ID).First(ctx)

	response := mappers.ToTournamentResponse(&result)
	return c.JSON(response)
}

func DeleteTournament(c *fiber.Ctx) error {
	ctx := c.Context()

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tournament ID",
		})
	}

	tournament, err := gorm.G[models.Tournament](db.DB).Where("id = ?", id).First(ctx)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Tournament not found",
		})
	}

	if _, err := gorm.G[models.Tournament](db.DB).Where("id = ?", tournament.ID).Delete(ctx); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete tournament",
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
