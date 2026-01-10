package routes

import (
	"github.com/PI-Team04-GameClub/gameclub-backend/handlers"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupTournamentRoutes(api fiber.Router, db *gorm.DB) {
	tournamentHandler := handlers.NewTournamentHandler(db)
	api.Get("/tournaments", tournamentHandler.GetTournaments)
	api.Get("/tournaments/:id", tournamentHandler.GetTournamentByID)
	api.Post("/tournaments", tournamentHandler.CreateTournament)
	api.Put("/tournaments/:id", tournamentHandler.UpdateTournament)
	api.Delete("/tournaments/:id", tournamentHandler.DeleteTournament)
}
