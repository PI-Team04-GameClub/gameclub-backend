package main

import (
	"log"

	"github.com/PI-Team04-GameClub/gameclub-backend/config"
	"github.com/PI-Team04-GameClub/gameclub-backend/db"
	"github.com/PI-Team04-GameClub/gameclub-backend/middleware"
	"github.com/PI-Team04-GameClub/gameclub-backend/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cfg := config.GetFromEnv()

	db.Connect(cfg)
	db.Migrate()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(logger.New())
	app.Use(middleware.SecurityHeaders())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))

	routes.Setup(app, db.DB)

	log.Fatal(app.Listen(":3000"))
}
