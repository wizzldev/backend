package events

import (
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/pkg/logger"
	"github.com/wizzldev/chat/pkg/ws"
)

func DispatchUserJoin(wsID string, userIDs []uint, user *models.User, groupID uint) error {
	m := &models.Message{
		HasGroup:         models.HasGroupID(groupID),
		HasMessageSender: models.HasMessageSenderID(user.ID),
		Content:          "",
		Type:             "join",
		DataJSON:         "{}",
	}
	database.DB.Save(m)

	sentTo := ws.WebSocket.BroadcastToUsers(userIDs, wsID, ws.Message{
		Event:  "join",
		Data:   user,
		HookID: "#",
	})

	logger.WSSend(wsID, "message.join", user.ID, sentTo)

	return nil
}
