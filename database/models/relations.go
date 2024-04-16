package models

type HasMessageSender struct {
	SenderID uint `json:"-"`
	Sender   User `json:"sender" gorm:"constraint:OnDelete:CASCADE;foreignKey:SenderID"`
}

type HasGroup struct {
	GroupID uint  `json:"-"`
	Group   Group `json:"receiver" gorm:"constraint:OnDelete:CASCADE;foreignKey:GroupID"`
}

type HasUser struct {
	UserID uint `json:"-"`
	User   User `json:"sender" gorm:"constraint:OnDelete:CASCADE;foreignKey:UserID"`
}
