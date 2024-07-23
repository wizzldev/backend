package models

type Block struct {
	Base
	HasUser
	BlockedUserID uint `json:"-"`
	BlockedUser   User `json:"blocked" gorm:"constraint:OnDelete:CASCADE;foreignKey:BlockedUserID"`

	Reason string `json:"reason"`
}
