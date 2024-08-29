package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/pkg/repository"
)

type developers struct{}

var Developers developers

func (developers) GetApplications(c *fiber.Ctx) error {
	bots := repository.Bot.FindBotsForUserID(authUserID(c))
	return c.JSON(fiber.Map{
		"bots": bots,
	})
}
