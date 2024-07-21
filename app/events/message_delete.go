package events

import (
	"errors"
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/pkg/logger"
	"github.com/wizzldev/chat/pkg/repository"
	"github.com/wizzldev/chat/pkg/ws"
	"strconv"
)

func DispatchMessageDelete(wsID string, userIDs []uint, user *models.User, msg *ws.ClientMessage, deleteOthers bool) error {
	id, err := strconv.Atoi(msg.Content)
	if err != nil {
		return err
	}

	m := repository.Message.FindOne(uint(id))
	if m.Sender.ID != user.ID && !deleteOthers {
		return errors.New("cannot delete message")
	}
	m.Type = "deleted"
	m.Content = ""
	m.DataJSON = "{}"
	m.ReplyID = nil
	database.DB.Save(m)

	sentTo := ws.WebSocket.BroadcastToUsers(userIDs, wsID, ws.Message{
		Event:  "message.unSend",
		Data:   m.ID,
		HookID: msg.HookID,
	})

	logger.WSSend(wsID, "message.unSend", user.ID, sentTo)

	return nil
}
