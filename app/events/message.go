package events

import (
	"github.com/wizzldev/chat/database"
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

func DispatchMessage(wsID string, userIDs []uint, gID uint, user *models.User, msg *ws.ClientMessage) {
	message := models.Message{
		HasGroup: models.HasGroup{
			GroupID: gID,
		},
		HasMessageSender: models.HasMessageSender{
			SenderID: user.ID,
		},
		Content:  msg.Content,
		Type:     msg.Type,
		DataJSON: msg.DataJSON,
	}
	database.DB.Create(&message)

	ws.WebSocket[wsID].BroadcastToUsers(userIDs, ws.Message{
		Event: "message",
		Data: ChatMessage{
			MessageID: message.ID,
			Sender:    *user,
			Content:   message.Content,
			Type:      message.Type,
			DataJSON:  message.DataJSON,
		},
	})
}
