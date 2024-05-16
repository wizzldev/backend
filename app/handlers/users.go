package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/app/requests"
	"github.com/wizzldev/chat/pkg/repository"
)

type users struct{}

var Users users

func (users) FindByEmail(c *fiber.Ctx) error {
	data := validation[requests.Email](c)

	user := repository.User.FindByEmail(data.Email)

	if repository.User.IsBlocked(user.ID, authUserID(c)) || user.ID < 1 {
		return fiber.ErrNotFound
	}

	return c.JSON(user)
}
