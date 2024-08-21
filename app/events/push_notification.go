package events

import (
	"github.com/wizzldev/chat/pkg/repository"
	"github.com/wizzldev/chat/pkg/services"
	"github.com/wizzldev/chat/pkg/utils"
)

type PushNotificationData struct {
	Title string
	Body  string
	Image string
}

func DispatchPushNotification(userIDs []uint, gID uint, data PushNotificationData) error {
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

	return services.PushNotification.Send(tokens, gID, data.Title, data.Body, utils.GetAvatarURL(data.Image, 24))
}
