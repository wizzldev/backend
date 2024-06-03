package models

type Group struct {
	Base
	IsPrivateMessage bool          `json:"is_private_message,omitempty"`
	Users            []*User       `json:"members,omitempty" gorm:"constraint:OnDelete:CASCADE;many2many:group_user"`
	MemberRoles      []*MemberRole `json:"member_roles,omitempty" gorm:"constraint:OnDelete:CASCADE;many2many:group_role"`
	ImageURL         string        `json:"image_url,omitempty" gorm:"default:null"`
	Name             string        `json:"name,omitempty" gorm:"default:null"`
	Roles            []string      `json:"roles,omitempty" gorm:"default:'[]';serializer:json"`
}
