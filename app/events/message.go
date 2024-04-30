package events

import (
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/pkg/ws"
	"slices"
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

	sentTo := ws.WebSocket[wsID].BroadcastToUsers(userIDs, ws.Message{
		Event: "message",
		Data: ChatMessage{
			MessageID: message.ID,
			Sender:    *user,
			Content:   message.Content,
			Type:      message.Type,
			DataJSON:  message.DataJSON,
		},
	})

	if len(sentTo) < len(userIDs) {
		shouldFetchIDs := make([]uint, len(userIDs)-len(sentTo))
		for _, id := range userIDs {
			if !slices.Contains(sentTo, id) {
				shouldFetchIDs = append(shouldFetchIDs, id)
			}
		}
		go ShouldFetch(shouldFetchIDs, gID)
	}
}
