package events

import (
	"github.com/wizzldev/chat/pkg/ws"
)

func SendToGroup(gID string, userIDs []uint, message ws.Message) {
	ws.WebSocket.BroadcastToUsers(userIDs, gID, message)
}
