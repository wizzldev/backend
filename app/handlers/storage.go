package handlers

import "github.com/gofiber/fiber/v2"

type storage struct {
	BasePath string
}

var Storage = storage{}

func (storage) Init() error {
	return nil
}

func (storage) Get(c *fiber.Ctx) error {
	return nil
}
