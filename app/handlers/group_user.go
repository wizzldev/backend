package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/app/requests"
	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/pkg/repository"
)

type groupUser struct{}

var GroupUser groupUser

func (g groupUser) EditNickName(c *fiber.Ctx) error {
	gu, err := g.helpData(c)
	if err != nil {
		return err
	}

	data := validation[requests.Nickname](c)
	gu.NickName = &data.Nickname

	err = repository.GroupUser.Update(gu)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"status": "ok",
	})
}

func (g groupUser) RemoveNickName(c *fiber.Ctx) error {
	gu, err := g.helpData(c)
	if err != nil {
		return err
	}

	gu.NickName = nil
	err = repository.GroupUser.Update(gu)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"status": "ok",
	})
}

func (groupUser) helpData(c *fiber.Ctx) (*models.GroupUser, error) {
	gID, err := c.ParamsInt("id")
	if err != nil {
		return nil, err
	}

	userID, err := c.ParamsInt("userID")
	if err != nil {
		return nil, err
	}

	gu, err := repository.GroupUser.Find(uint(gID), uint(userID))

	return gu, err
}
