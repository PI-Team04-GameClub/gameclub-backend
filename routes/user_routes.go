package routes

import (
	"github.com/PI-Team04-GameClub/gameclub-backend/handlers"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	usersBasePath = "/users"
	usersByIDPath = usersBasePath + "/:id"
)

func SetupUserRoutes(api fiber.Router, db *gorm.DB) {
	userHandler := handlers.NewUserHandler(db)
	api.Get(usersBasePath, userHandler.GetAllUsers)
	api.Get(usersByIDPath, userHandler.GetUserByID)
	api.Put(usersByIDPath, userHandler.UpdateUser)
}
