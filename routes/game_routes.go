package routes

import (
	"github.com/PI-Team04-GameClub/gameclub-backend/config"
	"github.com/PI-Team04-GameClub/gameclub-backend/handlers"
	"github.com/PI-Team04-GameClub/gameclub-backend/redis"
	"github.com/PI-Team04-GameClub/gameclub-backend/repositories"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	gamesBasePath = "/games"
	gamesByIDPath = gamesBasePath + "/:id"
)

func SetupGameRoutes(api fiber.Router, db *gorm.DB, cfg *config.Config) {
	baseRepo := repositories.NewGameRepository(db)

	var gameHandler *handlers.GameHandler
	if redis.Client != nil {
		redisCache := redis.NewRedisCache(redis.Client)
		cachedRepo := repositories.NewCachedGameRepository(baseRepo, redisCache, cfg.CacheTTL)
		gameHandler = handlers.NewGameHandlerWithRepo(cachedRepo)
	} else {
		gameHandler = handlers.NewGameHandlerWithRepo(baseRepo)
	}

	api.Get(gamesBasePath, gameHandler.GetAllGames)
	api.Get(gamesByIDPath, gameHandler.GetGameByID)
	api.Post(gamesBasePath, gameHandler.CreateGame)
	api.Put(gamesByIDPath, gameHandler.UpdateGame)
	api.Delete(gamesByIDPath, gameHandler.DeleteGame)
}
