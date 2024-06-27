package models

type AllowedIP struct {
	Base
	HasUser
	IP           string `json:"ip"`
	Active       bool   `json:"-"`
	Verification string `json:"-"`
}
