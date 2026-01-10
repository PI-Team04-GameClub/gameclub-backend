package interfaces

import "github.com/gofiber/fiber/v2"

type ITeamHandler interface {
	GetAllTeams(c *fiber.Ctx) error
	GetTeamByID(c *fiber.Ctx) error
	CreateTeam(c *fiber.Ctx) error
	UpdateTeam(c *fiber.Ctx) error
	DeleteTeam(c *fiber.Ctx) error
}
