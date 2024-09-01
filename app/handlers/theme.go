package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/pkg/repository"
)

type theme struct{}

var Theme theme

func (theme) Paginate(c *fiber.Ctx) error {
	data, err := repository.Theme.Paginate(c.Query("cursor"))
	if err != nil {
		return err
	}

	return c.JSON(data)
}
