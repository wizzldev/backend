package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/pkg/repository"
	"github.com/wizzldev/chat/pkg/utils"
)

func Auth(c *fiber.Ctx) error {
	sess, err := Session(c)
	if err != nil {
		return err
	}
	defer sess.Save()

	userId := sess.Get(utils.SessionAuthUserID)
	if userId == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	user := repository.User.FindById(userId.(uint))

	if user.ID < 1 {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
	}

	c.Locals(utils.LocalAuthUser, user)
	c.Locals(utils.LocalAuthUserID, user.ID)
	return c.Next()
}

func WSAuth(c *fiber.Ctx) error {
	q := c.Query("authorization", "none")
	if q != "none" {
		c.Request().Header.Set("Authorization", "bearer "+q)
	}
	return Auth(c)
}
