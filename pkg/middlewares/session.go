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

func AuthSession(c *fiber.Ctx) (*session.Session, error) {
	token, err := getToken(c.Request().Header.Peek("Authorization"))
	if err != nil {
		return nil, err
	}
	c.Request().Header.Set("Authorization", token)
	return Store.Get(c)
}

func Session(c *fiber.Ctx) (*session.Session, error) {
	return Store.Get(c)
}

func BotToken(c *fiber.Ctx) (string, error) {
	token, err := getToken(c.Request().Header.Peek("Authorization"))
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(token, "bot ") {
		return "", fiber.NewError(fiber.StatusBadRequest, "Authorization header does not contain bot token")
	}
	token = strings.TrimPrefix(token, "bot ")
	c.Request().Header.Set("Authorization", token)
	return token, nil
}

func getToken(raw []byte) (string, error) {
	authHeader := strings.ToLower(string(raw))
	if !strings.HasPrefix(authHeader, "bearer ") && authHeader != "" {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization header")
	}
	authHeaderTrimmed := strings.TrimPrefix(authHeader, "bearer ")
	if !strings.HasPrefix(authHeader, "bearer ") {
		return "", fiber.NewError(fiber.StatusForbidden, "Authorization header is not bearer token")
	}
	return authHeaderTrimmed, nil
}
