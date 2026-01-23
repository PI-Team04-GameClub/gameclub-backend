package routes

import (
	"github.com/PI-Team04-GameClub/gameclub-backend/handlers"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupTeamRoutes(api fiber.Router, db *gorm.DB) {
	teamHandler := handlers.NewTeamHandler(db)
	api.Get("/teams", teamHandler.GetAllTeams)
	api.Get("/teams/:id", teamHandler.GetTeamByID)
	api.Post("/teams", teamHandler.CreateTeam)
	api.Put("/teams/:id", teamHandler.UpdateTeam)
	api.Delete("/teams/:id", teamHandler.DeleteTeam)
	api.Get("/teams/:id/members", teamHandler.GetTeamMembers)
	api.Post("/teams/:id/members/:userId", teamHandler.JoinTeam)
}
