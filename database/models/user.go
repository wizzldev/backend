package models

import (
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
