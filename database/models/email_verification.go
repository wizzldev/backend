package models

type EmailVerification struct {
	Base
	HasUser
	Token string `json:"token" gorm:"token"`
}
