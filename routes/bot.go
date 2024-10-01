package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/pkg/middlewares"
)

func RegisterBot(r fiber.Router) {
	r.Get("/auth", middlewares.BotAuth, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
		})
	}).Name("auth")
}
