package repository

import (
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/pkg/repository/paginator"
)

type message struct{}

var Message message

func (message) FindOne(messageID uint) *models.Message {
	var msg models.Message

	_ = database.DB.Model(&models.Message{}).
		Preload("Sender").
		Preload("Reply.Sender").
		Preload("MessageLikes.User").
		Where("id = ?", messageID).
		Order("created_at desc").
		First(&msg).Error

	return &msg
}

func (message) Latest(gID uint) (*[]models.Message, string, string) {
	var messages []models.Message

	_ = database.DB.Model(&models.Message{}).
		Preload("Sender").
		Preload("Reply.Sender").
		Preload("MessageLikes.User").
		Where("group_id = ?", gID).
		Order("created_at desc").
		Limit(30).Find(&messages).Error

	return &messages, "", ""
}

func (message) CursorPaginate(gID uint, cursor string) (Pagination[models.Message], error) {
	query := database.DB.Model(&models.Message{}).Preload("Sender").
		Preload("Reply.Sender").
		Preload("MessageLikes.User").
		Where("group_id = ?", gID)

	data, next, prev, err := paginator.Paginate[models.Message](query, &paginator.Config{
		Cursor:     cursor,
		Order:      "desc",
		Limit:      30,
		PointsNext: false,
	})

	return Pagination[models.Message]{
		Data:       data,
		NextCursor: next,
		Previous:   prev,
	}, err
}
