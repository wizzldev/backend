package events

import (
	"github.com/wizzldev/chat/pkg/repository"
	"github.com/wizzldev/chat/pkg/services"
)

func DispatchPushNotification(userIDs []uint, gID uint, title, body, imageURL string) error {
	if len(userIDs) == 0 {
		return nil
	}

	err := services.PushNotification.Init()
	if err != nil {
		return err
	}

	tokens := repository.User.FindAndroidNotifications(userIDs)
	if len(tokens) == 0 {
		return nil
	}

	return services.PushNotification.Send(tokens, gID, title, body, imageURL)
}
