package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api")

	SetupAuthRoutes(api, db)
	SetupUserRoutes(api, db)
	SetupGameRoutes(api, db)
	SetupTeamRoutes(api, db)
	SetupTournamentRoutes(api, db)
	SetupNewsRoutes(api, db)
}
