package routes

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/app"
	"github.com/wizzldev/chat/app/handlers"
	"github.com/wizzldev/chat/pkg/middlewares"
	"github.com/wizzldev/chat/pkg/ws"
)

func WS(r fiber.Router) {
	ws.MessageHandler = app.WSActionHandler
	ws.Init()

	r.Use(middlewares.WSAuth)
	r.Use(func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	r.Get("/", handlers.WS.Connect)
	r.Use(HandleNotFoundError)
}
