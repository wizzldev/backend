package models

type ResetPassword struct {
	Base
	HasUser
	Token string `json:"token" gorm:"token"`
}
