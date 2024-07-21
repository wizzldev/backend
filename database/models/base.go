package models

import (
	"time"
)

type Base struct {
	ID        uint      `json:"id" gorm:"primaryKey;type:bigint(255)"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (b *Base) Exists() bool {
	return b.ID > 0
}
