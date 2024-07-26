package models

type AndroidPushNotification struct {
	Base
	HasUser
	Token string `json:"token"`
}
