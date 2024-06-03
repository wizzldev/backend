package ws

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/wizzldev/chat/pkg/configs"
	"github.com/wizzldev/chat/pkg/logger"
	"github.com/wizzldev/chat/pkg/utils"
	"slices"
	"sync"
)

var MessageHandler func(s *Server, conn *Connection, userID uint, message []byte) error

var WebSocket = map[string]*Server{}

type Server struct {
	ID   string
	Pool []*Connection
	mu   sync.Mutex
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
		Event:  "connection",
		Data:   "established",
		HookID: "#",
	})
	s.Pool = append(s.Pool, conn)

	if configs.Env.Debug {
		logger.WSNewConnection(s.ID, conn.IP(), conn.UserID)
	}

	s.LogPoolSize()

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
		if conn.Connected && f(conn) {
			conn.Send(m)
		}
	}
}

func (s *Server) BroadcastToUsers(userIDs []uint, m Message) []uint {
	var sentTo []uint
	for _, conn := range s.Pool {
		if slices.Contains(userIDs, conn.UserID) && conn.Connected {
			conn.Send(m)
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
	logger.WSPoolSize(s.ID, len(s.Pool), s.GetUserIDs())
}
