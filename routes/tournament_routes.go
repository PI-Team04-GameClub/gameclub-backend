package routes

import (
	"github.com/PI-Team04-GameClub/gameclub-backend/handlers"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	tournamentsBasePath = "/tournaments"
	tournamentsByIDPath = tournamentsBasePath + "/:id"
)

func SetupTournamentRoutes(api fiber.Router, db *gorm.DB) {
	tournamentHandler := handlers.NewTournamentHandler(db)
	api.Get(tournamentsBasePath, tournamentHandler.GetTournaments)
	api.Get(tournamentsByIDPath, tournamentHandler.GetTournamentByID)
	api.Post(tournamentsBasePath, tournamentHandler.CreateTournament)
	api.Put(tournamentsByIDPath, tournamentHandler.UpdateTournament)
	api.Delete(tournamentsByIDPath, tournamentHandler.DeleteTournament)
}
