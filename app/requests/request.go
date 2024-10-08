package requests

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/pkg/utils"
)

func Use[T any]() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return utils.Validate[T](c)
	}
}
