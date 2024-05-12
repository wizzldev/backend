package events

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/pkg/ws"
)

func ShouldFetch(userIDs []uint, gID uint) {
	fmt.Println("should fetch:", userIDs, gID)
	ws.Default().BroadcastToUsers(userIDs, ws.Message{
		Event: "should_fetch",
		Data: fiber.Map{
			"resource": fmt.Sprintf("chat.%v", gID),
			"group_id": gID,
		},
	})
}