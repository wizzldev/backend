package models

import (
	"gorm.io/datatypes"
)

type Group struct {
	Base
	IsPrivateMessage bool           `json:"is_private_message,omitempty"`
	Users            []*User        `json:"members" gorm:"constraint:OnDelete:CASCADE;many2many:group_user"`
	MemberRoles      []*MemberRole  `json:"member_roles" gorm:"constraint:OnDelete:CASCADE;many2many:group_role"`
	ImageURL         string         `json:"image_url,omitempty"`
	Name             string         `json:"name,omitempty"`
	Roles            datatypes.JSON `json:"roles"`
}
