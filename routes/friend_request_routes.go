package routes

import (
	"github.com/PI-Team04-GameClub/gameclub-backend/handlers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupFriendRequestRoutes(api fiber.Router, db *gorm.DB) {
	friendRequestHandler := handlers.NewFriendRequestHandler(db)

	friendRequests := api.Group("/friend-requests")
	friendRequests.Post("/", friendRequestHandler.CreateFriendRequest)
	friendRequests.Put("/:id/accept", friendRequestHandler.AcceptFriendRequest)
	friendRequests.Put("/:id/decline", friendRequestHandler.DeclineFriendRequest)
	friendRequests.Delete("/:id", friendRequestHandler.DeleteFriendRequest)

	// Get friend requests by user ID
	api.Get("/users/:id/friend-requests/sent", friendRequestHandler.GetSentFriendRequests)
	api.Get("/users/:id/friend-requests/received", friendRequestHandler.GetReceivedFriendRequests)

	// Get friends by user ID
	api.Get("/users/:id/friends", friendRequestHandler.GetFriends)
}
