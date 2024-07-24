package repository

import "time"

type LastMessage struct {
	SenderID   uint      `json:"sender_id"`
	SenderName string    `json:"sender_name"`
	NickName   *string   `json:"nick_name"`
	Content    *string   `json:"content"`
	Type       string    `json:"type"`
	Date       time.Time `json:"date"`
}

type Contact struct {
	ID               uint        `json:"id"`
	Name             string      `json:"name"`
	ImageURL         string      `json:"image"`
	Verified         bool        `json:"is_verified"`
	IsPrivateMessage bool        `json:"is_private_message"`
	CustomInvite     *string     `json:"custom_invite"`
	CreatorID        uint        `json:"creator_id"`
	LastMessage      LastMessage `json:"last_message"`
}
