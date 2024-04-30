package app

import (
	"github.com/wizzldev/chat/pkg/utils"
	"github.com/wizzldev/chat/pkg/ws"
)

func WSActionHandler(s *ws.Server, conn *ws.Connection, userID uint, data []byte) error {
	if s.ID != utils.DefaultWSPool {
		return MessageActionHandler(s, conn, userID, data)
	}
	return nil
}
