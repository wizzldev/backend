package models

type HasMessageSender struct {
	SenderID uint `json:"-"`
	Sender   User `json:"sender" gorm:"constraint:OnDelete:CASCADE;foreignKey:SenderID"`
}

func HasMessageSenderID(id uint) HasMessageSender {
	return HasMessageSender{SenderID: id}
}

type HasGroup struct {
	GroupID uint   `json:"-"`
	Group   *Group `json:"receiver,omitempty" gorm:"constraint:OnDelete:CASCADE;foreignKey:GroupID"`
}

func HasGroupID(id uint) HasGroup {
	return HasGroup{GroupID: id}
}

type HasUser struct {
	UserID uint `json:"-"`
	User   User `json:"sender" gorm:"constraint:OnDelete:CASCADE;foreignKey:UserID"`
}

func HasUserID(id uint) HasUser {
	return HasUser{UserID: id}
}

type HasMessageReply struct {
	ReplyID uint     `json:"-" gorm:"message_id;default:NULL"`
	Reply   *Message `json:"reply,omitempty" gorm:"constraint:OnDelete:CASCADE;foreignKey:ReplyID"`
}

type HasMessage struct {
	MessageID uint     `json:"-"`
	Message   *Message `json:"message,omitempty" gorm:"constraint:OnDelete:CASCADE;foreignKey:MessageID"`
}

type HasMessageLikes struct {
	MessageLikes []MessageLike `json:"likes,omitempty" gorm:"constraint:OnDelete:CASCADE;foreignKey:MessageID"`
}
