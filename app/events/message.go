package events

import (
	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/pkg/ws"
)

type ChatMessage struct {
	MessageID uint        `json:"message_id"`
	Sender    models.User `json:"sender"`
	Content   string      `json:"content"`
	Type      string      `json:"type"`
	DataJSON  string      `json:"data_json"`
}

func DispatchMessage(wsID string, userIDs []uint, m ChatMessage) {
	ws.WebSocket[wsID].BroadcastFunc(func(c *ws.Connection) bool {
		for _, id := range userIDs {
			if id == c.UserID {
				return true
			}
		}
		return false
	}, ws.Message{
		Event: "message.new",
		Data:  m,
	})
}
