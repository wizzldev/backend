package models

type Block struct {
	Base
	HasUser
	BlockedUserID uint `json:"-"`
	BlockedUser   User `json:"sender" gorm:"constraint:OnDelete:CASCADE;foreignKey:BlockedUserID"`

	Reason string `json:"reason"`
}
