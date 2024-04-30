package ws

import (
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"github.com/wizzldev/chat/pkg/configs"
	"github.com/wizzldev/chat/pkg/utils"
	"slices"
)

var MessageHandler func(s *Server, conn *Connection, userID uint, message []byte) error

var WebSocket = map[string]*Server{}

type Server struct {
	ID   string
	Pool []*Connection
}

type BroadcastFunc func(*Connection) bool

func Default() *Server {
	server, ok := WebSocket[utils.DefaultWSPool]

	if !ok {
		server = NewServer(utils.DefaultWSPool)
		WebSocket[utils.DefaultWSPool] = server
	}

	return server
}

func NewServer(id string) *Server {
	return &Server{
		ID:   id,
		Pool: []*Connection{},
	}
}

func (s *Server) AddConnection(ws *websocket.Conn) {
	conn := NewConnection(s.ID, ws, ws.Locals(utils.LocalAuthUserID).(uint))
	defer conn.Disconnect()
	conn.Init()
	conn.Send(Message{
		"connection",
		"established",
	})
	s.Pool = append(s.Pool, conn)

	if configs.Env.Debug {
		fmt.Printf("[WS] New connection, UserID: %v, IP: %v", conn.UserID, conn.IP())
	}

	conn.ReadLoop()
}

func (s *Server) Broadcast(m Message) {
	for _, conn := range s.Pool {
		if conn.Connected {
			conn.Send(m)
		}
	}
}

func (s *Server) BroadcastFunc(f BroadcastFunc, m Message) {
	for _, conn := range s.Pool {
		if f(conn) {
			conn.Send(m)
		}
	}
}

func (s *Server) BroadcastToUsers(userIDs []uint, m Message) []uint {
	var sentTo []uint
	for _, conn := range s.Pool {
		if slices.Contains(userIDs, conn.UserID) {
			conn.Send(m)
			sentTo = append(sentTo, conn.UserID)
		}
	}
	return sentTo
}

func (s *Server) Remove(c *Connection) {
	s.Pool = utils.RemoveFromSlice(s.Pool, c)
}
