package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/app/requests"
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/pkg/repository"
)

type mobile struct{}

var Mobile mobile

func (mobile) RegisterPushNotification(c *fiber.Ctx) error {
	userID := authUserID(c)
	data := validation[requests.PushToken](c)

	if !repository.User.IsAndroidNotificationTokenExists(userID, data.Token) {
		database.DB.Create(&models.AndroidPushNotification{
			HasUser: models.HasUserID(userID),
			Token:   data.Token,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
	})
}
