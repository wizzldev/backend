package models

import "gorm.io/datatypes"

type MemberRole struct {
	Base
	HasUser
	Roles datatypes.JSON `json:"roles"`
}
