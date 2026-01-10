package handlers

import (
	"strconv"

	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/mappers"
	"github.com/PI-Team04-GameClub/gameclub-backend/observer"
	"github.com/PI-Team04-GameClub/gameclub-backend/repositories"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TournamentHandler struct {
	tournamentRepo repositories.TournamentRepository
	gameRepo       repositories.GameRepository
	userRepo       repositories.UserRepository
}

func NewTournamentHandler(db *gorm.DB) *TournamentHandler {
	return &TournamentHandler{
		tournamentRepo: repositories.NewTournamentRepository(db),
		gameRepo:       repositories.NewGameRepository(db),
		userRepo:       repositories.NewUserRepository(db),
	}
}

func NewTournamentHandlerWithRepo(tournamentRepo repositories.TournamentRepository, gameRepo repositories.GameRepository, userRepo repositories.UserRepository) *TournamentHandler {
	return &TournamentHandler{
		tournamentRepo: tournamentRepo,
		gameRepo:       gameRepo,
		userRepo:       userRepo,
	}
}

func (h *TournamentHandler) GetTournaments(c *fiber.Ctx) error {
	tournaments, err := h.tournamentRepo.FindAll(c.Context())
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

	tournament, err := h.tournamentRepo.FindByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Tournament not found",
		})
	}

	response := mappers.ToTournamentResponse(tournament)
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

	_, err := h.gameRepo.FindByID(ctx, strconv.Itoa(int(req.GameId)))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Game not found",
		})
	}

	tournament := mappers.ToTournamentModel(req)

	if err := h.tournamentRepo.Create(ctx, &tournament); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create tournament",
		})
	}

	createdTournament, _ := h.tournamentRepo.FindByID(ctx, int(tournament.ID))

	users, _ := h.userRepo.FindAll(ctx)

	userEmails := make(map[string]string)
	for _, user := range users {
		userEmails[user.Email] = user.FirstName
	}

	createdTournament.Attach(observer.NewEmailNotifier(userEmails))
	createdTournament.Attach(observer.NewLogNotifier())

	createdTournament.NotifyCreated()

	response := mappers.ToTournamentResponse(createdTournament)
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

	tournament, err := h.tournamentRepo.FindByID(ctx, id)
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

	_, err = h.gameRepo.FindByID(ctx, strconv.Itoa(int(req.GameId)))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Game not found",
		})
	}

	updatedTournament := mappers.UpdateTournamentFromRequest(tournament, req)

	if err := h.tournamentRepo.Update(ctx, updatedTournament); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update tournament",
		})
	}

	result, _ := h.tournamentRepo.FindByID(ctx, int(updatedTournament.ID))

	response := mappers.ToTournamentResponse(result)
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

	_, err = h.tournamentRepo.FindByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Tournament not found",
		})
	}

	if err := h.tournamentRepo.Delete(ctx, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete tournament",
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
