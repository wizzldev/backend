package tests //nolint:typecheck

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/pkg/configs"
	"path/filepath"
)

func NewApp(envPath string) *fiber.App {
	err := configs.LoadEnv(filepath.Join(envPath, ".env.test"))
	if err != nil {
		log.Fatal(err)
	}

	return fiber.New(fiber.Config{
		Prefork:            !configs.Env.Debug,
		ServerHeader:       "Wizzl",
		AppName:            "Wizzl v1.0.0",
		ProxyHeader:        fiber.HeaderXForwardedFor,
		EnableIPValidation: true,
	})
}

func CleanUp() error {
	return database.CleanUpTestDB()
}
