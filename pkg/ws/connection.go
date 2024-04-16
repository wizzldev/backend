package ws

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2/log"
	"github.com/wizzldev/chat/pkg/configs"
)

type Connection struct {
	*websocket.Conn
	Connected bool
	UserID    uint
	serverID  string
}

func NewConnection(serverID string, ws *websocket.Conn, userId uint) *Connection {
	return &Connection{
		Conn:      ws,
		Connected: true,
		UserID:    userId,
		serverID:  serverID,
	}
}

func (c *Connection) Init() {
	return
}

func (c *Connection) Disconnect() {
	c.Connected = false
	_ = c.Conn.Close()
}

func (c *Connection) ReadLoop() {
	var (
		mt  int
		msg []byte
		err error
	)
	for c.Connected {
		if mt, msg, err = c.ReadMessage(); err != nil {
			c.Disconnect()
			break
		}
		if mt != websocket.TextMessage {
			continue
		}

		if configs.Env.Debug {
			c.Send(Message{
				Event: "echo",
				Data:  string(msg),
			})
		}

		err := MessageHandler(WebSocket[c.serverID], c, c.UserID, msg)

		if err != nil {
			if configs.Env.Debug {
				log.Warn("WS Read error:", err)
			}
			continue
		}
	}
}

func (c *Connection) Send(m Message) {
	if !c.Connected {
		return
	}
	err := c.Conn.WriteJSON(m)
	if err != nil {
		c.Disconnect()
	}
}
