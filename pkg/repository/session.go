package repository

import (
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
)

type session struct{}

var Session session

func (session) AllForUser(uID uint) []*models.Session {
	var sessions []*models.Session
	database.DB.Model(&models.Session{}).Where("user_id = ?", uID).Order("created_at desc").Limit(30).Find(&sessions)
	return sessions
}

func (session) FindForUser(uID uint, id uint) *models.Session {
	var s *models.Session
	database.DB.Model(&models.Session{}).Where("user_id = ? and id = ?", uID, id).First(&s)
	return s
}
