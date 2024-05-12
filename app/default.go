package app

import (
	"fmt"
	"github.com/wizzldev/chat/database/rdb"
	"github.com/wizzldev/chat/pkg/configs"
	"github.com/wizzldev/chat/pkg/utils"
	"github.com/wizzldev/chat/pkg/ws"
	"time"
)

func WSActionHandler(s *ws.Server, conn *ws.Connection, userID uint, data []byte) error {
	if configs.Env.Debug {
		fmt.Printf("WS[%v] New event: %s by user ID %v\n", s.ID, string(data), userID)
	}

	_ = rdb.Redis.Set(fmt.Sprintf("user.is-online.%v", userID), []byte("true"), time.Minute*10)
	if s.ID != utils.DefaultWSPool {
		return MessageActionHandler(s, conn, userID, data)
	}

	msg, err := ws.NewClientMessage(data, conn)

	if err != nil {
		return err
	}

	if msg.Type == "ping" {
		conn.Send(ws.Message{
			Event: "pong",
			Data:  "pong",
		})
	}

	return nil
}
