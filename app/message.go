package app

import (
	"fmt"
	"github.com/wizzldev/chat/app/events"
	"github.com/wizzldev/chat/app/services"
	"github.com/wizzldev/chat/pkg/ws"
	"slices"
	"strconv"
)

var cache = services.NewWSCache()

func MessageActionHandler(conn *ws.Connection, userID uint, msg *ws.ClientMessage, id string) error {
	user, err := cache.GetUser(userID)

	if err != nil {
		go conn.Send(ws.MessageWrapper{
			Message: &ws.Message{
				Event: "error",
				Data:  err.Error(),
			},
			Resource: id,
		})
		return err
	}

	gID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	members := cache.GetGroupMemberIDs(id)
	if !slices.Contains(members, userID) {
		return fmt.Errorf("user %d not in group members", userID)
	}

	switch msg.Type {
	case "message":
		return events.DispatchMessage(id, members, uint(gID), user, msg)
	case "message.like":
		return events.DispatchMessageLike(id, members, uint(gID), user, msg)
	case "message.delete":
		return events.DispatchMessageDelete(id, members, uint(gID), user, msg)
	default:
		conn.Send(ws.MessageWrapper{
			Message: &ws.Message{
				Event: "error",
				Data:  fmt.Sprintf("Unknown message type: %s", msg.Type),
			},
			Resource: id,
		})
	}

	return nil
}
