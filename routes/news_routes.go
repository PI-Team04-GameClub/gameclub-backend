package routes

import (
	"github.com/PI-Team04-GameClub/gameclub-backend/handlers"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	newsBasePath = "/news"
	newsByIDPath = newsBasePath + "/:id"
)

func SetupNewsRoutes(api fiber.Router, db *gorm.DB) {
	newsHandler := handlers.NewNewsHandler(db)
	api.Get(newsBasePath, newsHandler.GetNews)
	api.Post(newsBasePath, newsHandler.CreateNews)
	api.Put(newsByIDPath, newsHandler.UpdateNews)
	api.Delete(newsByIDPath, newsHandler.DeleteNews)
}
