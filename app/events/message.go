package events

import (
	"encoding/json"
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/pkg/logger"
	"github.com/wizzldev/chat/pkg/repository"
	"github.com/wizzldev/chat/pkg/utils"
	"github.com/wizzldev/chat/pkg/ws"
	"strings"
	"time"
)

type ChatMessage struct {
	MessageID uint            `json:"id"`
	Sender    models.User     `json:"sender"`
	Content   string          `json:"content"`
	Type      string          `json:"type"`
	DataJSON  string          `json:"data_json"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Reply     *models.Message `json:"reply"`
}

type DataJSON struct {
	ReplyID uint `json:"reply_id"`
}

func DispatchMessage(wsID string, userIDs []uint, gID uint, user *models.User, msg *ws.ClientMessage) error {
	var dataJSON DataJSON
	err := json.Unmarshal([]byte(msg.DataJSON), &dataJSON)
	if err != nil {
		return err
	}

	message := &models.Message{
		HasGroup: models.HasGroup{
			GroupID: gID,
		},
		HasMessageSender: models.HasMessageSender{
			SenderID: user.ID,
		},
		Content:  strings.TrimSpace(msg.Content),
		Type:     msg.Type,
		DataJSON: msg.DataJSON,
	}
	if dataJSON.ReplyID > 0 {
		message.ReplyID = &dataJSON.ReplyID
		message.DataJSON = "{}"
	}

	database.DB.Create(message)
	message = repository.Message.FindOne(message.ID)

	if user.IsBot {
		userIDs = utils.RemoveFromSlice(userIDs, user.ID)
	}

	sentTo := ws.WebSocket.BroadcastToUsers(userIDs, wsID, ws.Message{
		Event: "message",
		Data: ChatMessage{
			MessageID: message.ID,
			Sender:    *user,
			Content:   message.Content,
			Type:      message.Type,
			DataJSON:  message.DataJSON,
			CreatedAt: message.CreatedAt,
			UpdatedAt: message.UpdatedAt,
			Reply:     message.Reply,
		},
		HookID: msg.HookID,
	})

	logger.WSSend(wsID, "message", user.ID, sentTo)
	return nil
}
