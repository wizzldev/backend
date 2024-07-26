package database

import "github.com/wizzldev/chat/database/models"

func getModels() []interface{} {
	return []interface{}{
		&models.Message{},
		&models.MessageLike{},
		&models.Worker{},
		&models.AndroidPushNotification{},
		&models.Group{},
		&models.Ban{},
		&models.Invite{},
		&models.Block{},
		&models.EmailVerification{},
		&models.ResetPassword{},
		&models.Theme{},
		&models.GroupUser{},
		&models.AllowedIP{},
		&models.Session{},
		&models.User{},
		&models.File{},
	}
}
