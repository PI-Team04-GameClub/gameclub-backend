package interfaces

import "github.com/gofiber/fiber/v2"

type IGameHandler interface {
	GetAllGames(c *fiber.Ctx) error
	GetGameByID(c *fiber.Ctx) error
	CreateGame(c *fiber.Ctx) error
	UpdateGame(c *fiber.Ctx) error
	DeleteGame(c *fiber.Ctx) error
}
