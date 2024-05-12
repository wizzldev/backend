package repository

import "time"

type LastMessage struct {
	SenderID   uint      `json:"sender_id"`
	SenderName string    `json:"sender_name"`
	Content    *string   `json:"content"`
	Type       string    `json:"type"`
	Date       time.Time `json:"date"`
}

type Contact struct {
	ID          uint        `json:"id"`
	Name        string      `json:"name"`
	ImageURL    string      `json:"image"`
	LastMessage LastMessage `json:"last_message"`
}