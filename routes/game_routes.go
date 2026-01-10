package routes

import (
	"github.com/PI-Team04-GameClub/gameclub-backend/handlers"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupGameRoutes(api fiber.Router, db *gorm.DB) {
	gameHandler := handlers.NewGameHandler(db)
	api.Get("/games", gameHandler.GetAllGames)
	api.Get("/games/:id", gameHandler.GetGameByID)
	api.Post("/games", gameHandler.CreateGame)
	api.Put("/games/:id", gameHandler.UpdateGame)
	api.Delete("/games/:id", gameHandler.DeleteGame)
}
