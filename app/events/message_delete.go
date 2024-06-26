package events

import (
	"errors"
	"fmt"
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/pkg/logger"
	"github.com/wizzldev/chat/pkg/repository"
	"github.com/wizzldev/chat/pkg/ws"
	"strconv"
)

func DispatchMessageDelete(wsID string, userIDs []uint, gID uint, user *models.User, msg *ws.ClientMessage) error {
	id, err := strconv.Atoi(msg.Content)
	if err != nil {
		return err
	}

	m := repository.Message.FindOne(uint(id))
	if m.Sender.ID != user.ID {
		return errors.New("cannot delete message")
	}
	m.Type = "deleted"
	m.Content = ""
	m.DataJSON = "{}"
	fmt.Println("xxxxxxxxxxxxxxxxx", m.ReplyID, "data xxxxxxxx", m.Reply)
	database.DB.Save(m)

	sentTo := ws.WebSocket[wsID].BroadcastToUsers(userIDs, ws.Message{
		Event:  "message.unSend",
		Data:   m.ID,
		HookID: msg.HookID,
	})

	logger.WSSend(wsID, "message.unSend", user.ID, sentTo)

	DispatchShouldFetch(sentTo, userIDs, gID)
	return nil
}
