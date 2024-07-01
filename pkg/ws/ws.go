package ws

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/wizzldev/chat/pkg/configs"
	"github.com/wizzldev/chat/pkg/logger"
	"github.com/wizzldev/chat/pkg/utils"
	"slices"
	"sync"
)

var MessageHandler func(conn *Connection, userID uint, message []byte) error

var WebSocket *Server

type Server struct {
	Pool []*Connection
	mu   sync.Mutex
}

type BroadcastFunc func(*Connection) bool

func Init() *Server {
	if WebSocket == nil {
		WebSocket = NewServer()
	}

	return WebSocket
}

func NewServer() *Server {
	return &Server{
		Pool: []*Connection{},
	}
}

func (s *Server) AddConnection(ws *websocket.Conn) {
	conn := NewConnection("", ws, ws.Locals(utils.LocalAuthUserID).(uint))
	defer conn.Disconnect()
	conn.Init()
	conn.Send(MessageWrapper{
		Message: &Message{
			Event:  "connection",
			Data:   "established",
			HookID: "#",
		},
		Resource: utils.DefaultWSResource,
	})
	s.Pool = append(s.Pool, conn)

	if configs.Env.Debug {
		logger.WSNewConnection("", conn.IP(), conn.UserID)
	}

	s.LogPoolSize()

	conn.ReadLoop()
}

func (s *Server) Broadcast(m MessageWrapper) {
	for _, conn := range s.Pool {
		if conn.Connected {
			conn.Send(m)
		}
	}
}

func (s *Server) BroadcastFunc(f BroadcastFunc, m MessageWrapper) {
	for _, conn := range s.Pool {
		if conn.Connected && f(conn) {
			conn.Send(m)
		}
	}
}

func (s *Server) BroadcastToUsers(userIDs []uint, id string, m Message) []uint {
	var sentTo []uint
	for _, conn := range s.Pool {
		if slices.Contains(userIDs, conn.UserID) && conn.Connected {
			conn.Send(MessageWrapper{
				Message:  &m,
				Resource: id,
			})
			sentTo = append(sentTo, conn.UserID)
		}
	}
	return sentTo
}

func (s *Server) Remove(c *Connection) {
	s.LogPoolSize()
	s.Pool = utils.RemoveFromSlice(s.Pool, c)
	s.LogPoolSize()
}

func (s *Server) GetUserIDs() []uint {
	var userIDs []uint
	for _, conn := range s.Pool {
		userIDs = append(userIDs, conn.UserID)
	}
	return userIDs
}

func (s *Server) LogPoolSize() {
	logger.WSPoolSize("", len(s.Pool), s.GetUserIDs())
}
