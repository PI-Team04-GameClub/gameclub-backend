package handlers

import (
	"strconv"

	"github.com/PI-Team04-GameClub/gameclub-backend/dtos"
	"github.com/PI-Team04-GameClub/gameclub-backend/mappers"
	"github.com/PI-Team04-GameClub/gameclub-backend/repositories"
	"github.com/PI-Team04-GameClub/gameclub-backend/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TeamHandler struct {
	teamRepo repositories.TeamRepository
	userRepo repositories.UserRepository
}

func NewTeamHandler(db *gorm.DB) *TeamHandler {
	return &TeamHandler{
		teamRepo: repositories.NewTeamRepository(db),
		userRepo: repositories.NewUserRepository(db),
	}
}

func NewTeamHandlerWithRepo(teamRepo repositories.TeamRepository, userRepo repositories.UserRepository) *TeamHandler {
	return &TeamHandler{teamRepo: teamRepo, userRepo: userRepo}
}

func (h *TeamHandler) GetAllTeams(c *fiber.Ctx) error {
	teams, err := h.teamRepo.FindAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.InternalServerError("Failed to fetch teams"))
	}

	return c.JSON(mappers.ToTeamResponseList(teams))
}

func (h *TeamHandler) GetTeamByID(c *fiber.Ctx) error {
	id := c.Params("id")

	team, err := h.teamRepo.FindByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.NotFound())
	}

	return c.JSON(mappers.ToTeamResponse(team))
}

func (h *TeamHandler) CreateTeam(c *fiber.Ctx) error {
	var req dtos.CreateTeamRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.BadRequest("Invalid request body"))
	}

	team := mappers.ToTeamModel(req)

	if err := h.teamRepo.Create(c.Context(), &team); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.InternalServerError("Failed to create team"))
	}

	return c.Status(fiber.StatusCreated).JSON(mappers.ToTeamResponse(&team))
}

func (h *TeamHandler) UpdateTeam(c *fiber.Ctx) error {
	id := c.Params("id")

	team, err := h.teamRepo.FindByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.NotFound())
	}

	var req dtos.CreateTeamRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.BadRequest("Invalid request body"))
	}

	team.Name = req.Name

	if err := h.teamRepo.Update(c.Context(), team); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.InternalServerError("Failed to update team"))
	}

	return c.JSON(mappers.ToTeamResponse(team))
}

func (h *TeamHandler) DeleteTeam(c *fiber.Ctx) error {
	id := c.Params("id")

	team, err := h.teamRepo.FindByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.NotFound())
	}

	if err := h.teamRepo.Delete(c.Context(), team.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.InternalServerError("Failed to delete team"))
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *TeamHandler) GetTeamMembers(c *fiber.Ctx) error {
	id := c.Params("id")

	team, err := h.teamRepo.FindByIDWithMembers(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.NotFound())
	}

	return c.JSON(mappers.ToUserResponseListFromPointers(team.Users))
}

func (h *TeamHandler) JoinTeam(c *fiber.Ctx) error {
	teamID := c.Params("id")
	userIDParam := c.Params("userId")

	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.BadRequest("Invalid user ID"))
	}

	team, err := h.teamRepo.FindByID(c.Context(), teamID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.NotFound())
	}

	user, err := h.userRepo.FindByID(c.Context(), uint(userID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.NotFound())
	}

	if err := h.teamRepo.AddMember(c.Context(), team, user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.InternalServerError("Failed to add member to team"))
	}

	return c.SendStatus(fiber.StatusNoContent)
}
