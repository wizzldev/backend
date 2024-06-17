package models

type Bot struct {
	Base
	HasUser
	Name     string
	ImageURL string `json:"image_url"`
	IsOnline bool   `json:"is_online" gorm:"-:all"`
	Token    string `json:"-"`
}
