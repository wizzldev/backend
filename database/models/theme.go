package models

type Theme struct {
	Base
	HasUser
	Name string `json:"name"`
	Data string `json:"data" gorm:"type:json"`
}
