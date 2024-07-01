package handlers

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/pkg/ws"
)

type wsHandler struct{}

var WS wsHandler

func (wsHandler) Connect(c *fiber.Ctx) error {
	return websocket.New(ws.Init().AddConnection)(c)
}
