package models

type Worker struct {
	Base
	Command string `json:"command" gorm:"type:varchar(100)"`
	Data    string `json:"data" gorm:"type:longtext"`
}
