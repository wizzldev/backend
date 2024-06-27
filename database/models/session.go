package models

type Session struct {
	Base
	HasUser
	IP        string `json:"ip_address"`
	SessionID string `json:"-"`
	Agent     string `json:"user_agent"`
	Current   bool   `json:"current" gorm:"-:all"`
}
