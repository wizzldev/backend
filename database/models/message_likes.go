package models

type MessageLike struct {
	Base
	HasMessage
	HasUser
	Emoji string `json:"emoji"`
}
