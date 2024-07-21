package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/pkg/configs"
	"github.com/wizzldev/chat/pkg/repository"
	"github.com/wizzldev/chat/pkg/utils/role"
)

func NewRoleMiddleware(r role.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		gID, err := c.ParamsInt("id")
		if err != nil {
			return err
		}

		userID := c.Locals(configs.LocalAuthUserID).(uint)
		g := repository.Group.Find(uint(gID))
		if g.ID < 1 {
			return fiber.ErrNotFound
		}

		roles := repository.Group.GetUserRoles(uint(gID), userID, *role.NewRoles(g.Roles))
		if !roles.Can(r) {
			return fiber.NewError(fiber.StatusForbidden, "You are not allowed to access this resource")
		}

		return c.Next()
	}
}
