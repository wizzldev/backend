package models

type GroupUser struct {
	NickName string   `gorm:"nick_name,default:NULL"`
	Roles    []string `json:"roles,omitempty" gorm:"default:null;serializer:json"`
}

func (GroupUser) TableName() string {
	return "group_user"
}
