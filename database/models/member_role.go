package models

type MemberRole struct {
	Base
	HasUser
	HasGroup
	Role string `json:"role"`
}
