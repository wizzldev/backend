package repository

import (
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
)

type session struct{}

var Session session

func (session) AllForUser(uID uint) *[]models.Session {
	var sessions []models.Session
	database.DB.Model(&models.Session{}).Where("user_id = ?", uID).Find(&sessions)
	return &sessions
}
