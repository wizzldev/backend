package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/pkg/configs"
)

func validation[T any](c *fiber.Ctx) *T {
	return c.Locals(configs.RequestValidation).(*T)
}

func authUser(c *fiber.Ctx) *models.User {
	return c.Locals(configs.LocalAuthUser).(*models.User)
}

func authUserID(c *fiber.Ctx) uint {
	return c.Locals(configs.LocalAuthUserID).(uint)
}
