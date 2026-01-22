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

func AccessForbidden() Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors["body"] = "access forbidden"
	return e
}

func NotFound() Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors["body"] = "resource not found"
	return e
}

func RequiresAuthentication() Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors["body"] = "requested content requires authentication"
	return e
}

func BadRequest(message string) Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors["body"] = message
	return e
}

func InternalServerError(message string) Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors["body"] = message
	return e
}

func Conflict(message string) Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors["body"] = message
	return e
}

func Unauthorized(message string) Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors["body"] = message
	return e
}
