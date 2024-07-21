package models

import (
	"time"
)

type Invite struct {
	Base
	HasUser
	HasGroup
	MaxUsage   *int       `json:"max_usage" gorm:"default:null"`
	Key        string     `json:"key" gorm:"index:unique"` // join.wizzl.app/invite_id -> wizzl.app/invite/invite_id
	Expiration *time.Time `json:"expiration"`
}

// join.wizzl.app/releases
// join.wizzl.app/wizzl
// join.wizzl.app/support

func (i *Invite) IsValid() bool {
	// check if the invite is valid or not
	return i.ID > 0 &&
		(i.MaxUsage == nil || *i.MaxUsage > 0) &&
		(i.Expiration == nil || time.Now().Before(*i.Expiration))
}

func (i *Invite) Decrement() {
	if i.MaxUsage == nil {
		return
	}

	usage := *i.MaxUsage - 1
	i.MaxUsage = &usage
}
