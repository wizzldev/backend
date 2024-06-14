package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/app/handlers"
	"github.com/wizzldev/chat/pkg/middlewares"
)

func RegisterStorage(r fiber.Router) {
	file := r.Group("/files")
	file.Get("/:disc/:filename", middlewares.StorageFileToLocal(), middlewares.StorageFilePermission(), handlers.Files.Get)
	file.Use(HandleNotFoundError)

	avatar := r.Group("/avatars")
	avatar.Get("/:disc=s:size<int>.webp", middlewares.StorageFileToLocal(), handlers.Files.GetAvatar)
	avatar.Get("/:disc.webp", middlewares.StorageFileToLocal(), handlers.Files.GetAvatar)
	avatar.Use(HandleNotFoundError)

	r.Use(HandleNotFoundError)
}
