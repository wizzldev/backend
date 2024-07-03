package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/app/services"
	"net/http"
	"os"
	"time"
)

type files struct {
	*services.Storage
}

var Files = &files{}

func (s *files) Init(store *services.Storage) {
	s.Storage = store
}

func (s *files) Get(c *fiber.Ctx) error {
	return c.SendFile(s.LocalFile(c).Path)
}

func (s *files) GetInfo(c *fiber.Ctx) error {
	return c.JSON(s.LocalFile(c))
}

func (s *files) GetAvatar(c *fiber.Ctx) error {
	fileModel := s.LocalFile(c)
	if fileModel.Type != "avatar" {
		return fiber.ErrBadRequest
	}

	size, _ := c.ParamsInt("size")
	if size <= 0 {
		size = 256
	}

	file, err := os.Open(fileModel.Path)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to open file")
	}

	stream, err := s.WebPStream(file, uint(size))
	if err != nil {
		return err
	}

	c.Set("Content-Type", "image/webp")
	c.Set("Content-Disposition", "inline; filename=\"avatar.webp\"")
	c.Set("Cache-Control", "public, max-age=3600")
	c.Set("Last-Modified", fileModel.UpdatedAt.Format(http.TimeFormat))
	c.Set("Expires", time.Now().Add(24*time.Hour).Format(http.TimeFormat))
	return c.SendStream(stream)
}
