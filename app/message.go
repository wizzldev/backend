package app

import (
	"fmt"
	"github.com/wizzldev/chat/app/events"
	"github.com/wizzldev/chat/app/services"
	"github.com/wizzldev/chat/pkg/ws"
	"strconv"
)

var cache = services.NewWSCache()

func MessageActionHandler(s *ws.Server, conn *ws.Connection, userID uint, msg *ws.ClientMessage) error {
	user, err := cache.GetUser(userID)

	if err != nil {
		go conn.Send(ws.Message{
			Event: "error",
			Data:  err.Error(),
		})
		return err
	}

	gID, err := strconv.Atoi(s.ID)
	if err != nil {
		return err
	}

	members := cache.GetGroupMemberIDs(s.ID)

	switch msg.Type {
	case "message":
		return events.DispatchMessage(s.ID, members, uint(gID), user, msg)
	case "message.like":
		return events.DispatchMessageLike(s.ID, members, uint(gID), user, msg)
	case "message.delete":
		return events.DispatchMessageDelete(s.ID, members, uint(gID), user, msg)
	default:
		conn.Send(ws.Message{
			Event: "error",
			Data:  fmt.Sprintf("Unknown message type: %s", msg.Type),
		})
	}

	return nil
}
