package middleware

import "github.com/gofiber/fiber/v2"

// SecurityHeaders adds security headers to all responses
// Fixes:
// - X-Content-Type-Options Header Missing (CWE-693)
// - Strict-Transport-Security Header Not Set (CWE-319)
func SecurityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Prevent MIME-sniffing attacks
		c.Set("X-Content-Type-Options", "nosniff")

		// Enforce HTTPS connections (HSTS)
		// max-age=31536000 = 1 year
		// includeSubDomains = apply to all subdomains
		c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Additional recommended security headers
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")

		return c.Next()
	}
}
