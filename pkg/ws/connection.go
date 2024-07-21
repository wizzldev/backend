package ws

import (
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2/log"
	"github.com/wizzldev/chat/pkg/configs"
	"io"
	"sync"
	"time"
)

type Connection struct {
	*websocket.Conn
	Connected bool
	UserID    uint
	serverID  string
	mu        *sync.RWMutex
}

func NewConnection(serverID string, ws *websocket.Conn, userId uint) *Connection {
	return &Connection{
		Conn:      ws,
		Connected: true,
		UserID:    userId,
		serverID:  serverID,
		mu:        &sync.RWMutex{},
	}
}

func (c *Connection) Disconnect(msg ...string) {
	if !c.Connected {
		return
	}

	closeMessage := "closed by client"
	if len(msg) > 0 {
		closeMessage = msg[0]
	}

	WebSocket.mu.Lock()
	_ = c.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, closeMessage), time.Now().Add(time.Second))
	c.Connected = false
	err := c.Conn.Close()
	if err != nil {
		fmt.Println("Error closing connection", err)
	}
	WebSocket.Remove(c)
	fmt.Printf("Disconnected from server %s: %s \n", c.serverID, c.IP())
	WebSocket.mu.Unlock()
}

func (c *Connection) ReadLoop() {
	var (
		mt  int
		msg []byte
		err error
	)
	for c.Connected {
		if mt, msg, err = c.ReadMessage(); err != nil {
			break
		}
		if mt != websocket.TextMessage {
			continue
		}

		if configs.Env.Debug {
			c.Send(MessageWrapper{
				Message: &Message{
					Event:  "echo",
					Data:   string(msg),
					HookID: "#",
				},
				Resource: configs.DefaultWSResource,
			})
		}

		err := MessageHandler(c, c.UserID, msg)

		if err != nil {
			if configs.Env.Debug {
				log.Warn("WS Read error:", err)
			}

			if err != io.EOF {
				c.Send(MessageWrapper{
					Message: &Message{
						Event:  "error",
						Data:   err.Error(),
						HookID: "#",
					},
					Resource: "#",
				})
			}

			continue
		}
	}
	c.Disconnect()
}

func (c *Connection) Send(m MessageWrapper) {
	if !c.Connected {
		return
	}
	c.mu.Lock()
	err := c.Conn.WriteJSON(m)
	c.mu.Unlock()
	if err != nil {
		c.Disconnect()
	}
}
