package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/pkg/configs"
	"github.com/wizzldev/chat/pkg/repository"
	"strings"
)

func Auth(c *fiber.Ctx) error {
	sess, err := AuthSession(c)
	if err != nil {
		return err
	}

	userId := sess.Get(configs.SessionAuthUserID)
	if userId == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	user := repository.User.FindById(userId.(uint))

	if user.ID < 1 {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	if user.IsBot {
		return fiber.NewError(fiber.StatusBadRequest, "Cannot use bot as a user")
	}

	c.Locals(configs.LocalAuthUser, user)
	c.Locals(configs.LocalAuthUserID, user.ID)
	c.Locals(configs.LocalIsBot, user.IsBot)
	return c.Next()
}

func WSAuth(c *fiber.Ctx) error {
	q := c.Query("authorization", "none")

	if q != "none" {
		c.Request().Header.Set("Authorization", "bearer "+q)
	}

	return AnyAuth(c)
}

func BotAuth(c *fiber.Ctx) error {
	token, err := BotToken(c)
	if err != nil {
		return err
	}
	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	bot := repository.Bot.FindByToken(token)

	if bot.ID < 1 {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	c.Locals(configs.LocalAuthUser, bot)
	c.Locals(configs.LocalAuthUserID, bot.ID)
	c.Locals(configs.LocalIsBot, bot.IsBot)
	return c.Next()
}

func AnyAuth(c *fiber.Ctx) error {
	if strings.Contains(strings.ToLower(string(c.Request().Header.Peek("Authorization"))), " bot ") {
		return BotAuth(c)
	}
	return Auth(c)
}

func NoBots(c *fiber.Ctx) error {
	if c.Locals(configs.LocalIsBot).(bool) {
		return fiber.NewError(fiber.StatusForbidden, "Bots not allowed")
	}
	return c.Next()
}
