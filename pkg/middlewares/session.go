package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/pkg/configs"
	utils2 "github.com/wizzldev/chat/pkg/utils"
	"strings"
	"time"
)

var store = session.New(session.Config{
	Expiration: time.Duration(configs.Env.Session.LifespanSeconds) * time.Second,
	KeyGenerator: func() string {
		return "bearer " + utils.UUIDv4()
	},
	KeyLookup: "header:Authorization",
	Storage:   database.Redis,
})

func Session(c *fiber.Ctx) (*session.Session, error) {
	authHeader := strings.ToLower(string(c.Request().Header.Peek("Authorization")))
	if !strings.HasPrefix(authHeader, "bearer ") || !utils2.IsValidUUID(strings.TrimPrefix(authHeader, "bearer ")) {
		c.Request().Header.Del("Authorization")
	}
	c.Request().Header.Set("Authorization", authHeader)
	return store.Get(c)
}
