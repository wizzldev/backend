package models

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/database/rdb"
	"gorm.io/gorm"
	"time"
)

type User struct {
	Base
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	Email           string     `json:"-"`
	Password        string     `json:"-"`
	ImageURL        string     `json:"image_url"`
	EmailVerifiedAt *time.Time `json:"-"`
	IsOnline        bool       `json:"is_online" gorm:"-:all"`
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

func (u *User) AfterFind(*gorm.Tx) error {
	exists, _ := rdb.RedisClient.Exists(ctx, fmt.Sprintf("user.is-online.%v", u.ID)).Result()
	u.IsOnline = exists == 1
	fmt.Println("after find:", u.ID, u.IsOnline)
	return nil
}
