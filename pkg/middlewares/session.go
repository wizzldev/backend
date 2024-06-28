package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/wizzldev/chat/database/rdb"
	"github.com/wizzldev/chat/pkg/configs"
	utils2 "github.com/wizzldev/chat/pkg/utils"
	"strings"
	"time"
)

var Store = session.New(session.Config{
	Expiration: time.Duration(configs.Env.Session.LifespanSeconds) * time.Second,
	KeyGenerator: func() string {
		return "client:" + utils.UUIDv4() + "$" + strings.ToLower(utils2.NewRandom().String(25))
	},
	KeyLookup: "header:Authorization",
	Storage:   rdb.Redis,
})

func Session(c *fiber.Ctx) (*session.Session, error) {
	authHeader := strings.ToLower(string(c.Request().Header.Peek("Authorization")))
	authHeaderTrimmed := strings.TrimPrefix(authHeader, "bearer ")
	if !strings.HasPrefix(authHeader, "bearer ") {
		c.Request().Header.Del("Authorization")
	}
	c.Request().Header.Set("Authorization", authHeaderTrimmed)
	return Store.Get(c)
}
