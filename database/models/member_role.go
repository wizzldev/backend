package models

type MemberRole struct {
	Base
	HasUser
	Roles string `json:"roles"`
}
