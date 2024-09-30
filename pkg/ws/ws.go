package ws

import (
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/wizzldev/chat/pkg/configs"
	"github.com/wizzldev/chat/pkg/logger"
	"github.com/wizzldev/chat/pkg/utils"
)

var MessageHandler func(conn *Connection, userID uint, message []byte) error

var WebSocket *Server

type Server struct {
	Pool []*Connection
	mu   sync.RWMutex
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
		Pool: make([]*Connection, 0, 100),
	}
}

func (s *Server) AddConnection(ws *websocket.Conn) {
	conn := NewConnection("", ws, ws.Locals(configs.LocalAuthUserID).(uint))
	defer conn.Disconnect()

	conn.Send(MessageWrapper{
		Message: &Message{
			Event:  "connection",
			Data:   "established",
			HookID: "#",
		},
		Resource: configs.DefaultWSResource,
	})

	s.mu.Lock()
	s.Pool = append(s.Pool, conn)
	s.mu.Unlock()

	if configs.Env.Debug {
		logger.WSNewConnection("", conn.IP(), conn.UserID)
	}

	s.LogPoolSize()

	conn.ReadLoop()
}

func (s *Server) Broadcast(m MessageWrapper) {
	s.mu.RLock() // Read lock since we aren't modifying the pool
	defer s.mu.RUnlock()

	var wg sync.WaitGroup
	for _, conn := range s.Pool {
		if conn.Connected {
			wg.Add(1)
			go func(c *Connection) {
				defer wg.Done()
				c.Send(m)
			}(conn)
		}
	}
	wg.Wait() // Wait for all goroutines to finish sending
}

func (s *Server) BroadcastFunc(f BroadcastFunc, m MessageWrapper) {
	s.mu.RLock() // Read lock
	defer s.mu.RUnlock()

	var wg sync.WaitGroup
	for _, conn := range s.Pool {
		if conn.Connected && f(conn) {
			wg.Add(1)
			go func(c *Connection) {
				defer wg.Done()
				c.Send(m)
			}(conn)
		}
	}
	wg.Wait() // Wait for all goroutines to finish
}

func (s *Server) BroadcastToUsers(userIDs []uint, id string, m Message) []uint {
	userIDMap := make(map[uint]struct{}, len(userIDs))
	for _, id := range userIDs {
		userIDMap[id] = struct{}{}
	}

	var sentTo []uint
	s.mu.RLock() // Read lock
	defer s.mu.RUnlock()

	var wg sync.WaitGroup
	mu := sync.Mutex{} // Mutex to safely append to sentTo
	for _, conn := range s.Pool {
		if _, exists := userIDMap[conn.UserID]; exists && conn.Connected {
			wg.Add(1)
			go func(c *Connection) {
				defer wg.Done()
				c.Send(MessageWrapper{
					Message:  &m,
					Resource: id,
				})
				mu.Lock()
				sentTo = append(sentTo, c.UserID)
				mu.Unlock()
			}(conn)
		}
	}
	wg.Wait()
	return sentTo
}

func (s *Server) Remove(c *Connection) {
	s.mu.Lock()
	s.Pool = utils.RemoveFromSlice(s.Pool, c)
	s.mu.Unlock()

	s.LogPoolSize()
}

func (s *Server) GetUserIDs() []uint {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userIDs := make([]uint, 0, len(s.Pool))
	for _, conn := range s.Pool {
		userIDs = append(userIDs, conn.UserID)
	}
	return userIDs
}

func (s *Server) LogPoolSize() {
	userIDs := s.GetUserIDs()
	logger.WSPoolSize("", len(s.Pool), userIDs)
}
