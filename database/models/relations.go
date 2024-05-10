package models

type HasMessageSender struct {
	SenderID uint `json:"-"`
	Sender   User `json:"sender" gorm:"constraint:OnDelete:CASCADE;foreignKey:SenderID"`
}

type HasGroup struct {
	GroupID uint   `json:"-"`
	Group   *Group `json:"receiver,omitempty" gorm:"constraint:OnDelete:CASCADE;foreignKey:GroupID"`
}

type HasUser struct {
	UserID uint `json:"-"`
	User   User `json:"sender" gorm:"constraint:OnDelete:CASCADE;foreignKey:UserID"`
}

type GroupUser struct {
	GroupID uint `gorm:"group_id"`
	UserID  uint `gorm:"user_id"`
}

type HasMessageReply struct {
	ReplyID uint     `json:"-" gorm:"message_id;default:NULL"`
	Reply   *Message `json:"reply,omitempty" gorm:"constraint:OnDelete:CASCADE;foreignKey:ReplyID"`
}
