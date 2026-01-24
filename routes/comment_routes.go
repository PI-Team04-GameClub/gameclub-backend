package routes

import (
	"github.com/PI-Team04-GameClub/gameclub-backend/handlers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupCommentRoutes(api fiber.Router, db *gorm.DB) {
	commentHandler := handlers.NewCommentHandler(db)

	comments := api.Group("/comments")
	comments.Post("/", commentHandler.CreateComment)
	comments.Put("/:id", commentHandler.UpdateComment)

	api.Get("/users/:id/comments", commentHandler.GetCommentsByUserID)

	api.Get("/news/:id/comments", commentHandler.GetCommentsByNewsID)
}
