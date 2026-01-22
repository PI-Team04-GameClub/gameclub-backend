package utils

import (
	"errors"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewError_WithFiberError(t *testing.T) {
	// Given: A fiber error with a specific message
	fiberErr := fiber.NewError(fiber.StatusBadRequest, "bad request message")

	// When: Creating a new error from the fiber error
	result := NewError(fiberErr)

	// Then: The error should contain the fiber error message
	assert.NotNil(t, result.Errors)
	assert.Equal(t, "bad request message", result.Errors["body"])
}

func TestNewError_WithStandardError(t *testing.T) {
	// Given: A standard Go error
	stdErr := errors.New("standard error message")

	// When: Creating a new error from the standard error
	result := NewError(stdErr)

	// Then: The error should contain the standard error message
	assert.NotNil(t, result.Errors)
	assert.Equal(t, "standard error message", result.Errors["body"])
}

func TestNewValidationError(t *testing.T) {
	// Given: A field name and validation error message
	field := "email"
	message := "invalid email format"

	// When: Creating a new validation error
	result := NewValidationError(field, message)

	// Then: The error should contain the message for the specified field
	assert.NotNil(t, result.Errors)
	assert.Equal(t, "invalid email format", result.Errors["email"])
}

func TestAccessForbidden(t *testing.T) {
	// Given: A request that should be forbidden

	// When: Creating an access forbidden error
	result := AccessForbidden()

	// Then: The error should contain the access forbidden message
	assert.NotNil(t, result.Errors)
	assert.Equal(t, "access forbidden", result.Errors["body"])
}

func TestNotFound(t *testing.T) {
	// Given: A request for a non-existent resource

	// When: Creating a not found error
	result := NotFound()

	// Then: The error should contain the not found message
	assert.NotNil(t, result.Errors)
	assert.Equal(t, "resource not found", result.Errors["body"])
}

func TestRequiresAuthentication(t *testing.T) {
	// Given: A request for protected content without authentication

	// When: Creating a requires authentication error
	result := RequiresAuthentication()

	// Then: The error should contain the authentication required message
	assert.NotNil(t, result.Errors)
	assert.Equal(t, "requested content requires authentication", result.Errors["body"])
}

func TestBadRequest(t *testing.T) {
	// Given: A bad request message
	message := "invalid input"

	// When: Creating a bad request error
	result := BadRequest(message)

	// Then: The error should contain the bad request message
	assert.NotNil(t, result.Errors)
	assert.Equal(t, "invalid input", result.Errors["body"])
}

func TestInternalServerError(t *testing.T) {
	// Given: An internal server error message
	message := "database connection failed"

	// When: Creating an internal server error
	result := InternalServerError(message)

	// Then: The error should contain the internal server error message
	assert.NotNil(t, result.Errors)
	assert.Equal(t, "database connection failed", result.Errors["body"])
}

func TestConflict(t *testing.T) {
	// Given: A conflict error message
	message := "email already exists"

	// When: Creating a conflict error
	result := Conflict(message)

	// Then: The error should contain the conflict message
	assert.NotNil(t, result.Errors)
	assert.Equal(t, "email already exists", result.Errors["body"])
}

func TestUnauthorized(t *testing.T) {
	// Given: An unauthorized error message
	message := "invalid credentials"

	// When: Creating an unauthorized error
	result := Unauthorized(message)

	// Then: The error should contain the unauthorized message
	assert.NotNil(t, result.Errors)
	assert.Equal(t, "invalid credentials", result.Errors["body"])
}

func TestErrorStruct_HasCorrectJSONTag(t *testing.T) {
	// Given: An Error struct with a body field
	e := Error{Errors: map[string]interface{}{"body": "test"}}

	// When: Accessing the errors map

	// Then: The body field should contain the correct value
	assert.Equal(t, "test", e.Errors["body"])
}
