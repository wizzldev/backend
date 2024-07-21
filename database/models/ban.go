package models

type Ban struct {
	Base
	HasGroup
	HasUser

	BlockedUserID uint `json:"-"`
	BlockedUser   User `json:"sender" gorm:"constraint:OnDelete:CASCADE;foreignKey:BlockedUserID"`

	Reason string `json:"reason"`
}
