package events

import (
	"github.com/wizzldev/chat/pkg/push_notification"
	"github.com/wizzldev/chat/pkg/repository"
)

func DispatchPushNotification(userIDs []uint, title, body, imageURL string) error {
	if len(userIDs) == 0 {
		return nil
	}

	err := push_notification.PushNotification.Init()
	if err != nil {
		return err
	}

	tokens := repository.User.FindAndroidNotifications(userIDs)
	if len(tokens) == 0 {
		return nil
	}

	return push_notification.PushNotification.Send(tokens, title, body, imageURL)
}
