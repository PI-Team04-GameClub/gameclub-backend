package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestSecurityHeaders(t *testing.T) {
	// Given:  A Fiber app with security headers middleware
	app := fiber.New()
	app.Use(SecurityHeaders())

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// When: Making a request
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	// Then: Security headers should be present
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Verify X-Frame-Options header
	assert.Equal(t, "DENY", resp.Header.Get("X-Frame-Options"))

	// Verify Content-Security-Policy header
	assert.Equal(t, "frame-ancestors 'none'", resp.Header.Get("Content-Security-Policy"))

	// Verify additional security headers
	assert.Equal(t, "nosniff", resp.Header.Get("X-Content-Type-Options"))
	assert.Equal(t, "1; mode=block", resp.Header.Get("X-XSS-Protection"))
	assert.Equal(t, "strict-origin-when-cross-origin", resp.Header.Get("Referrer-Policy"))
}

func TestSecurityHeaders_AppliedToAllRoutes(t *testing.T) {
	// Given: A Fiber app with multiple routes
	app := fiber.New()
	app.Use(SecurityHeaders())

	app.Get("/route1", func(c *fiber.Ctx) error {
		return c.SendString("Route 1")
	})

	app.Post("/route2", func(c *fiber.Ctx) error {
		return c.SendString("Route 2")
	})

	// When: Making requests to different routes
	tests := []struct {
		method string
		path   string
	}{
		{"GET", "/route1"},
		{"POST", "/route2"},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.path, nil)
		resp, err := app.Test(req)

		// Then: All routes should have the security headers
		assert.NoError(t, err)
		assert.Equal(t, "DENY", resp.Header.Get("X-Frame-Options"))
	}
}
