package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/app/handlers"
	"github.com/wizzldev/chat/pkg/middlewares"
)

func RegisterStorage(r fiber.Router) {
	file := r.Group("/files").Name("files.")
	file.Get("/:disc/:filename", middlewares.StorageFileToLocal(), middlewares.StorageFilePermission(), handlers.Files.Get).Name("file")
	file.Get("/:disc/:filename/info", middlewares.StorageFileToLocal(), middlewares.StorageFilePermission(), handlers.Files.GetInfo).Name("info")
	file.Use(HandleNotFoundError)

	avatar := r.Group("/avatars").Name("avatars.")
	avatar.Post("/upload", func(c *fiber.Ctx) error {
		fh, err := c.FormFile("image")
		if err != nil {
			return err
		}
		f, err := handlers.Files.StoreAvatar(fh)
		if err != nil {
			return err
		}
		return c.SendString(f.Discriminator)
	}).Name("upload")
	avatar.Get("/:disc-s:size<int>.webp", middlewares.StorageFileToLocal(), handlers.Files.GetAvatar).Name("with-size")
	avatar.Get("/:disc.webp", middlewares.StorageFileToLocal(), handlers.Files.GetAvatar).Name("simple")
	avatar.Use(HandleNotFoundError)

	r.Use(HandleNotFoundError)
}
