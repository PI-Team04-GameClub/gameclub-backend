package routes

import (
	"github.com/PI-Team04-GameClub/gameclub-backend/handlers"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	gamesBasePath = "/games"
	gamesByIDPath = gamesBasePath + "/:id"
)

func SetupGameRoutes(api fiber.Router, db *gorm.DB) {
	gameHandler := handlers.NewGameHandler(db)
	api.Get(gamesBasePath, gameHandler.GetAllGames)
	api.Get(gamesByIDPath, gameHandler.GetGameByID)
	api.Post(gamesBasePath, gameHandler.CreateGame)
	api.Put(gamesByIDPath, gameHandler.UpdateGame)
	api.Delete(gamesByIDPath, gameHandler.DeleteGame)
}
