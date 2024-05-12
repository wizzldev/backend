package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/app/handlers"
)

func RegisterStorage(r fiber.Router) {
	if err := handlers.Storage.Init(); err != nil {
		panic(err)
	}

	r.Get("/:resource/:file", handlers.Storage.Get)
}
