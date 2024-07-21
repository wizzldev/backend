package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/wizzldev/chat/app/requests"
	"github.com/wizzldev/chat/pkg/configs"
	"time"
)

func NewSimpleLimiter(max int, expiration time.Duration, message string, keyGenerator ...func(*fiber.Ctx) string) fiber.Handler {
	var gen = func(c *fiber.Ctx) string {
		return c.IP()
	}

	if len(keyGenerator) == 1 {
		gen = keyGenerator[0]
	}

	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: expiration,
		LimitReached: func(c *fiber.Ctx) error {
			return fiber.NewError(fiber.StatusTooManyRequests, message)
		},
		KeyGenerator: gen,
	})
}

func NewAuthLimiter() fiber.Handler {
	return NewSimpleLimiter(10, 10*time.Minute, "Too many attempts, try again later", func(c *fiber.Ctx) string {
		req := c.Locals(configs.RequestValidation).(*requests.Login)
		return req.Email
	})
}
