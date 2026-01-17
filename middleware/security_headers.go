package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// SecurityHeaders middleware adds security-related headers to all responses
func SecurityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// X-Frame-Options:  Prevents clickjacking attacks
		c.Set("X-Frame-Options", "DENY") 
		// DENY - Prevents any domain from framing the content

		// Content-Security-Policy: Modern alternative to X-Frame-Options
		// frame-ancestors 'none' is equivalent to X-Frame-Options: DENY
		c.Set("Content-Security-Policy", "frame-ancestors 'none'")

		// X-Content-Type-Options:  Prevents MIME type sniffing
		c.Set("X-Content-Type-Options", "nosniff")

		// X-XSS-Protection: Enables browser's XSS filter
		c.Set("X-XSS-Protection", "1; mode=block")

		// Referrer-Policy: Controls referrer information
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")

		return c.Next()
	}
}
