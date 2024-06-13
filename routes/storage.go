package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/app/handlers"
	"github.com/wizzldev/chat/pkg/middlewares"
)

func RegisterStorage(r fiber.Router) {
	if err := handlers.Storage.Init(); err != nil {
		panic(err)
	}

	file := r.Group("/files/:disc")
	file.Get("/:filename", middlewares.StorageFileToLocal(), middlewares.StorageFilePermission(), handlers.Storage.Get)
	file.Use(HandleNotFoundError)

	avatar := r.Group("/avatars")
	avatar.Post("/upload", middlewares.Auth, handlers.Storage.StoreAvatar)
	avatar.Get("/:disc.webp", middlewares.StorageFileToLocal(), handlers.Storage.GetAvatar)
	avatar.Use(HandleNotFoundError)

	r.Use(HandleNotFoundError)
}
