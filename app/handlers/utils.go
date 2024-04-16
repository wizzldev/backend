package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/pkg/utils"
)

func validation[T any](c *fiber.Ctx) *T {
	return c.Locals("requestValidation").(*T)
}

func authUser(c *fiber.Ctx) *models.User {
	return c.Locals(utils.LocalAuthUser).(*models.User)
}

func authUserId(c *fiber.Ctx) uint {
	return c.Locals(utils.LocalAuthUserID).(uint)
}
