package utils

import (
	"errors"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewError_WithFiberError(t *testing.T) {
	fiberErr := fiber.NewError(fiber.StatusBadRequest, "bad request message")
	result := NewError(fiberErr)

	assert.NotNil(t, result.Errors)
	assert.Equal(t, "bad request message", result.Errors["body"])
}

func TestNewError_WithStandardError(t *testing.T) {
	stdErr := errors.New("standard error message")
	result := NewError(stdErr)

	assert.NotNil(t, result.Errors)
	assert.Equal(t, "standard error message", result.Errors["body"])
}

func TestNewValidationError(t *testing.T) {
	result := NewValidationError("email", "invalid email format")

	assert.NotNil(t, result.Errors)
	assert.Equal(t, "invalid email format", result.Errors["email"])
}

func TestAccessForbidden(t *testing.T) {
	result := AccessForbidden()

	assert.NotNil(t, result.Errors)
	assert.Equal(t, "access forbidden", result.Errors["body"])
}

func TestNotFound(t *testing.T) {
	result := NotFound()

	assert.NotNil(t, result.Errors)
	assert.Equal(t, "resource not found", result.Errors["body"])
}

func TestRequiresAuthentication(t *testing.T) {
	result := RequiresAuthentication()

	assert.NotNil(t, result.Errors)
	assert.Equal(t, "requested content requires authentication", result.Errors["body"])
}

func TestBadRequest(t *testing.T) {
	result := BadRequest("invalid input")

	assert.NotNil(t, result.Errors)
	assert.Equal(t, "invalid input", result.Errors["body"])
}

func TestInternalServerError(t *testing.T) {
	result := InternalServerError("database connection failed")

	assert.NotNil(t, result.Errors)
	assert.Equal(t, "database connection failed", result.Errors["body"])
}

func TestConflict(t *testing.T) {
	result := Conflict("email already exists")

	assert.NotNil(t, result.Errors)
	assert.Equal(t, "email already exists", result.Errors["body"])
}

func TestUnauthorized(t *testing.T) {
	result := Unauthorized("invalid credentials")

	assert.NotNil(t, result.Errors)
	assert.Equal(t, "invalid credentials", result.Errors["body"])
}

func TestErrorStruct_HasCorrectJSONTag(t *testing.T) {
	e := Error{Errors: map[string]interface{}{"body": "test"}}

	assert.Equal(t, "test", e.Errors["body"])
}
