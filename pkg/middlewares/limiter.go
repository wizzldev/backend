package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/wizzldev/chat/app/requests"
	"github.com/wizzldev/chat/pkg/utils"
	"time"
)

func NewAuthLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        10,
		Expiration: 5 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			req := c.Locals(utils.RequestValidation).(*requests.Login)
			return req.Email
		},
	})
}
