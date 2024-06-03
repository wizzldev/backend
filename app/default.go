package app

import (
	"fmt"
	"github.com/wizzldev/chat/database/rdb"
	"github.com/wizzldev/chat/pkg/configs"
	"github.com/wizzldev/chat/pkg/logger"
	"github.com/wizzldev/chat/pkg/utils"
	"github.com/wizzldev/chat/pkg/ws"
	"time"
)

func WSActionHandler(s *ws.Server, conn *ws.Connection, userID uint, data []byte) error {
	msg, err := ws.NewClientMessage(data, conn)

	if err != nil {
		return err
	}

	if configs.Env.Debug {
		logger.WSNewEvent(s.ID, msg.Type, userID)
	}

	if msg.Type == "ping" {
		conn.Send(ws.Message{
			Event:  "pong",
			Data:   "pong",
			HookID: msg.HookID,
		})
		return nil
	}

	if msg.Type == "close" {
		logger.WSDisconnect(s.ID, userID)
		conn.Disconnect()
		return nil
	}

	_ = rdb.Redis.Set(fmt.Sprintf("user.is-online.%v", userID), []byte("true"), time.Minute*10)
	if s.ID != utils.DefaultWSPool {
		return MessageActionHandler(s, conn, userID, msg)
	}

	return nil
}
