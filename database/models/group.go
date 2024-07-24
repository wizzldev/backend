package models

type Group struct {
	Base
	HasTheme
	HasUser
	IsPrivateMessage bool     `json:"is_private_message"`
	Users            []*User  `json:"members,omitempty" gorm:"constraint:OnDelete:CASCADE;many2many:group_user"`
	ImageURL         string   `json:"image_url,omitempty" gorm:"default:null"`
	Name             string   `json:"name,omitempty" gorm:"default:null"`
	Roles            []string `json:"roles,omitempty" gorm:"default:'[]';serializer:json"`
	Verified         bool     `json:"is_verified" gorm:"default:false"`
	Emoji            *string  `json:"emoji,omitempty" gorm:"default:null"`
	CustomInvite     *string  `json:"custom_invite,omitempty" gorm:"default:null"`
}
