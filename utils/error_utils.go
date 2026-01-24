package utils

import (
	"github.com/gofiber/fiber/v2"
)

type Error struct {
	Errors map[string]interface{} `json:"errors"`
}

func NewError(err error) Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	switch v := err.(type) {
	case *fiber.Error:
		e.Errors["body"] = v.Message
	default:
		e.Errors["body"] = v.Error()
	}
	return e
}

func NewValidationError(field string, message string) Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors[field] = message
	return e
}

func newBodyError(message string) Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors["body"] = message
	return e
}

func AccessForbidden() Error {
	return newBodyError("access forbidden")
}

func NotFound() Error {
	return newBodyError("resource not found")
}

func RequiresAuthentication() Error {
	return newBodyError("requested content requires authentication")
}

func BadRequest(message string) Error {
	return newBodyError(message)
}

func InternalServerError(message string) Error {
	return newBodyError(message)
}

func Conflict(message string) Error {
	return newBodyError(message)
}

func Unauthorized(message string) Error {
	return newBodyError(message)
}
