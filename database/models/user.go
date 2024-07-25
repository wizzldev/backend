package models

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/database/rdb"
	"gorm.io/gorm"
)

type User struct {
	Base
	FirstName       string     `json:"first_name" gorm:"type:varchar(100)"`
	LastName        string     `json:"last_name" gorm:"type:varchar(100)"`
	Email           string     `json:"email" gorm:"type:varchar(100)"`
	Password        string     `json:"-" gorm:"type:varchar(255)"`
	ImageURL        string     `json:"image_url" gorm:"type:varchar(255)"`
	EmailVerifiedAt *time.Time `json:"-"`
	EnableIPCheck   bool       `json:"enable_ip_check" gorm:"default:true"`
	IsOnline        bool       `json:"is_online" gorm:"-:all"`
	IsBot           bool       `json:"is_bot" gorm:"default:false"`
	GroupUser       *GroupUser `json:"group_user,omitempty"`
}

var ctx = context.Background()

func (u *User) PublicData() fiber.Map {
	return fiber.Map{
		"first_name": u.FirstName,
		"last_name":  u.LastName,
		"image_url":  u.ImageURL,
		"is_online":  u.IsOnline,
	}
}

func (u *User) AfterFind(*gorm.DB) error {
	exists, _ := rdb.RedisClient.Exists(ctx, fmt.Sprintf("user.is-online.%v", u.ID)).Result()
	u.IsOnline = exists == 1
	return nil
}
