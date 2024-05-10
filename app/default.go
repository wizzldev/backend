package app

import (
	"fmt"
	"github.com/wizzldev/chat/database/rdb"
	"github.com/wizzldev/chat/pkg/utils"
	"github.com/wizzldev/chat/pkg/ws"
	"time"
)

func WSActionHandler(s *ws.Server, conn *ws.Connection, userID uint, data []byte) error {
	_ = rdb.Redis.Set(fmt.Sprintf("user.is-online.%v", userID), []byte("true"), time.Minute*10)
	if s.ID != utils.DefaultWSPool {
		return MessageActionHandler(s, conn, userID, data)
	}
	return nil
}
