package events

import (
	"encoding/json"
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/pkg/repository"
	"github.com/wizzldev/chat/pkg/utils"
	"github.com/wizzldev/chat/pkg/ws"
	"strings"
)

type MessageLike struct {
	ID        uint         `json:"id"`
	Emoji     string       `json:"emoji"`
	Sender    *models.User `json:"sender"`
	MessageID uint         `json:"message_id"`
}

func DispatchMessageLike(wsID string, userIDs []uint, _ uint, user *models.User, msg *ws.ClientMessage) error {
	msgID := struct {
		MessageID uint `json:"message_id"`
	}{}

	err := json.NewDecoder(strings.NewReader(msg.DataJSON)).Decode(&msgID)
	if err != nil {
		return err
	}

	if repository.IsExists[models.Message]([]string{"id"}, []any{msgID.MessageID}) {
		messageLike := repository.FindModelBy[models.MessageLike]([]string{"message_id", "user_id"}, []any{msgID.MessageID, user.ID})
		var isLiked bool

		if messageLike.ID > 0 {
			database.DB.Delete(messageLike)
		} else {
			messageLike = &models.MessageLike{
				HasMessage: models.HasMessage{
					MessageID: msgID.MessageID,
				},
				HasUser: models.HasUser{
					UserID: user.ID,
				},
				Emoji: msg.Content,
			}
			database.DB.Create(messageLike)
			isLiked = true
		}

		var t = "like"
		if !isLiked {
			t += ".remove"
		}

		if user.IsBot {
			userIDs = utils.RemoveFromSlice(userIDs, user.ID)
		}

		_ = ws.WebSocket.BroadcastToUsers(userIDs, wsID, ws.Message{
			Event: "message." + t,
			Data: &MessageLike{
				ID:        messageLike.ID,
				Emoji:     messageLike.Emoji,
				Sender:    user,
				MessageID: messageLike.MessageID,
			},
			HookID: msg.HookID,
		})
	}

	return nil
}
