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

	gameHandler := handlers.NewGameHandler(db)
	api.Get("/games", gameHandler.GetAllGames)
	api.Get("/games/:id", gameHandler.GetGameByID)
	api.Post("/games", gameHandler.CreateGame)
	api.Put("/games/:id", gameHandler.UpdateGame)
	api.Delete("/games/:id", gameHandler.DeleteGame)

	teamHandler := handlers.NewTeamHandler(db)
	api.Get("/teams", teamHandler.GetAllTeams)
	api.Get("/teams/:id", teamHandler.GetTeamByID)
	api.Post("/teams", teamHandler.CreateTeam)
	api.Put("/teams/:id", teamHandler.UpdateTeam)
	api.Delete("/teams/:id", teamHandler.DeleteTeam)

	tournamentHandler := handlers.NewTournamentHandler(db)
	api.Get("/tournaments", tournamentHandler.GetTournaments)
	api.Get("/tournaments/:id", tournamentHandler.GetTournamentByID)
	api.Post("/tournaments", tournamentHandler.CreateTournament)
	api.Put("/tournaments/:id", tournamentHandler.UpdateTournament)
	api.Delete("/tournaments/:id", tournamentHandler.DeleteTournament)

	newsHandler := handlers.NewNewsHandler(db)
	api.Get("/news", newsHandler.GetNews)
	api.Post("/news", newsHandler.CreateNews)
	api.Put("/news/:id", newsHandler.UpdateNews)
	api.Delete("/news/:id", newsHandler.DeleteNews)
}
