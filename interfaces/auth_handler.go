package interfaces

import "github.com/gofiber/fiber/v2"

type IAuthHandler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	GetCurrentUser(c *fiber.Ctx) error
}
