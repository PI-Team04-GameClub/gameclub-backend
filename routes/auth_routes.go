package routes

import (
	"github.com/PI-Team04-GameClub/gameclub-backend/handlers"
	"github.com/PI-Team04-GameClub/gameclub-backend/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupAuthRoutes(api fiber.Router, db *gorm.DB) {
	authHandler := handlers.NewAuthHandler(db)
	api.Post("/auth/register", authHandler.Register)
	api.Post("/auth/login", authHandler.Login)
	api.Get("/auth/me", middleware.JWTMiddleware(db), authHandler.GetCurrentUser)
}
