package interfaces

import "github.com/gofiber/fiber/v2"

type ITournamentHandler interface {
	GetTournaments(c *fiber.Ctx) error
	GetTournamentByID(c *fiber.Ctx) error
	CreateTournament(c *fiber.Ctx) error
	UpdateTournament(c *fiber.Ctx) error
	DeleteTournament(c *fiber.Ctx) error
}
