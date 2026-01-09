package routes

import (
	"github.com/PI-Team04-GameClub/gameclub-backend/handlers"
	"github.com/PI-Team04-GameClub/gameclub-backend/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api")

	authHandler := handlers.NewAuthHandler(db)
	api.Post("/auth/register", authHandler.Register)
	api.Post("/auth/login", authHandler.Login)

	api.Get("/auth/me", middleware.JWTMiddleware(db), authHandler.GetCurrentUser)

	api.Get("/games", handlers.GetAllGames)
	api.Get("/games/:id", handlers.GetGameByID)
	api.Post("/games", handlers.CreateGame)
	api.Put("/games/:id", handlers.UpdateGame)
	api.Delete("/games/:id", handlers.DeleteGame)

	api.Get("/teams", handlers.GetAllTeams)
	api.Get("/teams/:id", handlers.GetTeamByID)
	api.Post("/teams", handlers.CreateTeam)
	api.Put("/teams/:id", handlers.UpdateTeam)
	api.Delete("/teams/:id", handlers.DeleteTeam)

	api.Get("/tournaments", handlers.GetTournaments)
	api.Get("/tournaments/:id", handlers.GetTournamentByID)
	api.Post("/tournaments", handlers.CreateTournament)
	api.Put("/tournaments/:id", handlers.UpdateTournament)
	api.Delete("/tournaments/:id", handlers.DeleteTournament)

	api.Get("/news", handlers.GetNews)
	api.Post("/news", handlers.CreateNews)
	api.Put("/news/:id", handlers.UpdateNews)
	api.Delete("/news/:id", handlers.DeleteNews)
}
