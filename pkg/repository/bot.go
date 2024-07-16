package repository

import (
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
)

type bot struct{}

var Bot bot

func (bot) FindByToken(t string) *models.User {
	var bot models.User
	database.DB.Model(&models.User{}).Where("password = ? and is_bot = ?", t, true).First(&bot)
	return &bot
}
