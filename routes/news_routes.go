package routes

import (
	"github.com/PI-Team04-GameClub/gameclub-backend/handlers"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupNewsRoutes(api fiber.Router, db *gorm.DB) {
	newsHandler := handlers.NewNewsHandler(db)
	api.Get("/news", newsHandler.GetNews)
	api.Post("/news", newsHandler.CreateNews)
	api.Put("/news/:id", newsHandler.UpdateNews)
	api.Delete("/news/:id", newsHandler.DeleteNews)
}
