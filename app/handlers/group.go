package handlers

import "github.com/gofiber/fiber/v2"

type group struct{}

var Group group

func (group) New(c *fiber.Ctx) error {
	return nil
}
