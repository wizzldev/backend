package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/app/services"
	"github.com/wizzldev/chat/database"
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
	user.ImageURL = fmt.Sprintf("{api}/storage/avatars/%s.webp", file.Discriminator)
	database.DB.Save(user)

	return c.JSON(user)
}
