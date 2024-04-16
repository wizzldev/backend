package routes

import (
	"errors"
	"github.com/gofiber/fiber/v2"
)

func HandleNotFoundError(*fiber.Ctx) error {
	return fiber.NewError(fiber.StatusNotFound, "This resource could not be found")
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
