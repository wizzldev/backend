package events

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/pkg/ws"
	"slices"
)

func ShouldFetch(userIDs []uint, gID uint) {
	ws.Default().BroadcastToUsers(userIDs, ws.Message{
		Event: "should_fetch",
		Data: fiber.Map{
			"resource": fmt.Sprintf("chat.%v", gID),
			"group_id": gID,
		},
	})
}

func DispatchShouldFetch(sentTo []uint, allUserID []uint, gID uint) {
	if len(sentTo) < len(allUserID) {
		var shouldFetchIDs []uint
		for _, id := range allUserID {
			if !slices.Contains(sentTo, id) {
				shouldFetchIDs = append(shouldFetchIDs, id)
			}
		}
		ShouldFetch(shouldFetchIDs, gID)
	}
}
