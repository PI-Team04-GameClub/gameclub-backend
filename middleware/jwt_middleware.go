package middleware

import (
	"context"
	"fmt"

	"github.com/PI-Team04-GameClub/gameclub-backend/models"
	"github.com/PI-Team04-GameClub/gameclub-backend/security"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func JWTMiddleware(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		claims, err := security.ExtractClaims(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		userID, ok := claims["id"].(float64)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid user ID in token",
			})
		}

		user, err := gorm.G[models.User](db).Where("id = ?", uint(userID)).First(context.Background())
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		c.Locals("user", &user)
		c.Locals("userID", uint(userID))

		return c.Next()
	}
}

func UnauthorizedHandler(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": fmt.Sprintf("Unauthorized: %v", err),
	})
}
