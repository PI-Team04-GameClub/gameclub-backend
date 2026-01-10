package interfaces

import "github.com/gofiber/fiber/v2"

type INewsHandler interface {
	GetNews(c *fiber.Ctx) error
	CreateNews(c *fiber.Ctx) error
	UpdateNews(c *fiber.Ctx) error
	DeleteNews(c *fiber.Ctx) error
}
