package requests

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/pkg/utils"
)

func Use(r any) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return utils.Validate(r, c)
	}
}
