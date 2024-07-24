package app

import (
	"errors"
	"fmt"
	"github.com/wizzldev/chat/app/events"
	"github.com/wizzldev/chat/app/services"
	"github.com/wizzldev/chat/pkg/utils/role"
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

	gIDInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	gID := uint(gIDInt)

	members := cache.GetGroupMemberIDs(id)
	if !slices.Contains(members, userID) {
		return fmt.Errorf("user %d not in group members", userID)
	}

	isPM := cache.IsPM(gID)
	roles := cache.GetRoles(userID, gID)
	roleErr := errors.New("you do not have a permit")

	switch msg.Type {
	case "message":
		if !roles.Can(role.SendMessage) && !isPM {
			return roleErr
		}
		return events.DispatchMessage(id, members, gID, user, msg)
	case "message.like":
		return events.DispatchMessageLike(id, members, gID, user, msg)
	case "message.delete":
		if !roles.Can(role.DeleteMessage) && !isPM {
			return roleErr
		}
		return events.DispatchMessageDelete(id, members, user, msg, roles.Can(role.DeleteOtherMemberMessage) && !isPM)
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
