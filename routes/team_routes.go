package routes

import (
	"github.com/PI-Team04-GameClub/gameclub-backend/handlers"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	teamsBasePath = "/teams"
	teamsByIDPath = teamsBasePath + "/:id"
)

func SetupTeamRoutes(api fiber.Router, db *gorm.DB) {
	teamHandler := handlers.NewTeamHandler(db)
	api.Get(teamsBasePath, teamHandler.GetAllTeams)
	api.Get(teamsByIDPath, teamHandler.GetTeamByID)
	api.Post(teamsBasePath, teamHandler.CreateTeam)
	api.Put(teamsByIDPath, teamHandler.UpdateTeam)
	api.Delete(teamsByIDPath, teamHandler.DeleteTeam)
	api.Get(teamsByIDPath+"/members", teamHandler.GetTeamMembers)
	api.Post(teamsByIDPath+"/members/:userId", teamHandler.JoinTeam)
}
