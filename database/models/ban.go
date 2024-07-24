package models

import "time"

type Ban struct {
	Base
	HasGroup
	HasUser

	BlockedUserID uint `json:"-"`
	BlockedUser   User `json:"sender" gorm:"constraint:OnDelete:CASCADE;foreignKey:BlockedUserID"`

	Duration *time.Time `json:"duration" gorm:"default:null"`

	Reason string `json:"reason"`
}
