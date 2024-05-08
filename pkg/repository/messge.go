package repository

import (
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
)

type message struct{}

var Message message

func (message) Latest(gID uint) *[]models.Message {
	var messages []models.Message

	_ = database.DB.Model(&models.Message{}).
		Preload("Sender").
		Where("group_id = ?", gID).
		Order("created_at desc").
		Limit(30).Find(&messages).Error

	return &messages
}
