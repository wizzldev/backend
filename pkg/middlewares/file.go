package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/pkg/utils"
	"net/url"
	"strings"
)

func StorageFileToLocal() fiber.Handler {
	return func(c *fiber.Ctx) error {
		disc := c.Params("disc")
		if strings.Contains(disc, "-") {
			s := strings.SplitN(disc, "-", 2)
			if len(s) == 2 {
				disc = s[0]
			}
		}

		var file models.File
		database.DB.Model(&models.File{}).
			Where("discriminator = ?", disc).
			Find(&file)

		if file.ID < 1 {
			return fiber.ErrNotFound
		}

		rawName := c.Params("filename")
		name, err := url.QueryUnescape(rawName)
		if err != nil {
			return err
		}

		if name != "" && name != file.Name {
			return fiber.ErrNotFound
		}

		c.Locals(utils.LocalFileModel, &file)
		return c.Next()
	}
}

func StorageFilePermission() fiber.Handler {
	return func(c *fiber.Ctx) error {
		file := c.Locals(utils.LocalFileModel).(*models.File)

		if file.AccessToken != nil && !canAccessFile(c, *file.AccessToken) {
			return fiber.ErrForbidden
		}

		return c.Next()
	}
}

func canAccessFile(c *fiber.Ctx, token string) bool {
	rawCode := c.Request().Header.Peek("X-File-Access-Token")
	return string(rawCode) == token || c.Query("access_token") == token
}
