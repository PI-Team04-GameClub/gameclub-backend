package routes

import (
	"github.com/PI-Team04-GameClub/gameclub-backend/handlers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	friendRequestsBasePath    = "/friend-requests"
	friendRequestsByIDPath    = "/:id"
	friendRequestUserBasePath = "/users"
	friendRequestUserByIDPath = friendRequestUserBasePath + "/:id"
)

func SetupFriendRequestRoutes(api fiber.Router, db *gorm.DB) {
	friendRequestHandler := handlers.NewFriendRequestHandler(db)

	friendRequests := api.Group(friendRequestsBasePath)
	friendRequests.Post("/", friendRequestHandler.CreateFriendRequest)
	friendRequests.Put(friendRequestsByIDPath+"/accept", friendRequestHandler.AcceptFriendRequest)
	friendRequests.Put(friendRequestsByIDPath+"/decline", friendRequestHandler.DeclineFriendRequest)
	friendRequests.Delete(friendRequestsByIDPath, friendRequestHandler.DeleteFriendRequest)

	// Get friend requests by user ID
	api.Get(friendRequestUserByIDPath+"/friend-requests/sent", friendRequestHandler.GetSentFriendRequests)
	api.Get(friendRequestUserByIDPath+"/friend-requests/received", friendRequestHandler.GetReceivedFriendRequests)

	// Get friends by user ID
	api.Get(friendRequestUserByIDPath+"/friends", friendRequestHandler.GetFriends)
}
