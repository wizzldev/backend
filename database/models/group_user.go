package models

type GroupUser struct {
	NickName string `gorm:"nick_name,default:NULL"`
}

func (GroupUser) TableName() string {
	return "group_user"
}
