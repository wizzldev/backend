package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/app/services"
	"github.com/wizzldev/chat/database"
	"strings"
)

type me struct {
	*services.Storage
}

var Me = &me{}

func (m *me) Init(store *services.Storage) {
	m.Storage = store
}

func (*me) Hello(c *fiber.Ctx) error {
	user := authUser(c)
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Hello %s", user.FirstName),
		"user":    user,
	})
}

func (m *me) UploadProfileImage(c *fiber.Ctx) error {
	fileH, err := c.FormFile("image")
	if err != nil {
		return err
	}

	file, err := m.Storage.StoreAvatar(fileH)
	if err != nil {
		return err
	}

	user := authUser(c)
	if user.ImageURL != "default.webp" {
		_ = m.Storage.RemoveByDisc(strings.SplitN(user.ImageURL, ".", 2)[0])
	}

	user.ImageURL = file.Discriminator + ".webp"
	database.DB.Save(user)

	return c.JSON(user)
}
