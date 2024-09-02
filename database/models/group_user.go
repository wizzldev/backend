package models

type GroupUser struct {
	HasGroup
	HasUser
	NickName *string  `json:"nick_name,omitempty" gorm:"default:NULL"`
	Roles    []string `json:"roles,omitempty" gorm:"default:null;serializer:json"`
}

func (GroupUser) TableName() string {
	return "group_user"
}
