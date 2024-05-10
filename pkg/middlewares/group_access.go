package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/pkg/utils"
)

func GroupAccess(IDLookup string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authUserID := c.Locals(utils.LocalAuthUserID).(uint)
		groupID, err := c.ParamsInt(IDLookup)
		if err != nil {
			return err
		}

		var can int64
		err = database.DB.Table("group_user").
			Where("group_id = ? and user_id = ?", groupID, authUserID).
			Limit(1).
			Count(&can).Error

		if err != nil || can < 1 {
			return fiber.ErrForbidden
		}

		return c.Next()
	}
}
