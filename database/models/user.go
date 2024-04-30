package models

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/database/rdb"
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
}

var ctx = context.Background()

func (u *User) PublicData() fiber.Map {
	err := rdb.RedisClient.Exists(ctx, fmt.Sprintf("user.is-online.%v", u.ID)).Err()
	isOnline := err == nil

	return fiber.Map{
		"first_name": u.FirstName,
		"last_name":  u.LastName,
		"image_url":  u.ImageURL,
		"is_online":  isOnline,
	}
}
