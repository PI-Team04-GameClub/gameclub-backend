package routes

import (
	"github.com/PI-Team04-GameClub/gameclub-backend/config"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB, cfg *config.Config) {
	api := app.Group("/api")

	SetupAuthRoutes(api, db)
	SetupUserRoutes(api, db)
	SetupGameRoutes(api, db, cfg)
	SetupTeamRoutes(api, db)
	SetupTournamentRoutes(api, db)
	SetupNewsRoutes(api, db)
	SetupCommentRoutes(api, db)
	SetupFriendRequestRoutes(api, db)
}
