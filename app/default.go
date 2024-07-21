package app

import (
	"fmt"
	"github.com/wizzldev/chat/database/rdb"
	"github.com/wizzldev/chat/pkg/configs"
	"github.com/wizzldev/chat/pkg/logger"
	"github.com/wizzldev/chat/pkg/ws"
	"time"
)

func WSActionHandler(conn *ws.Connection, userID uint, data []byte) error {
	wrapper, err := ws.NewMessage(data, conn)

	if err != nil {
		return err
	}

	msg := wrapper.Message

	if configs.Env.Debug {
		logger.WSNewEvent(wrapper.Resource, msg.Type, userID)
	}

	if msg.Type == "ping" {
		conn.Send(ws.MessageWrapper{
			Message: &ws.Message{
				Event:  "pong",
				Data:   "pong",
				HookID: msg.HookID,
			},
			Resource: configs.DefaultWSResource,
		})
		return nil
	}

	if msg.Type == "close" {
		logger.WSDisconnect(wrapper.Resource, userID)
		conn.Disconnect()
		return nil
	}

	_ = rdb.Redis.Set(fmt.Sprintf("user.is-online.%v", userID), []byte("true"), time.Minute*10)
	if wrapper.Resource != configs.DefaultWSResource {
		return MessageActionHandler(conn, userID, msg, wrapper.Resource)
	}

	return nil
}
