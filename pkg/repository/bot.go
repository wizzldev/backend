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

func (bot) FindBotsForUserID(userID uint) *[]models.User {
	var bots []models.User
	database.DB.Raw(`
		select users.* from users
		inner join user_bots on user_bots.bot_id = users.id
		where users.is_bot = 1 and user_bots.user_id = ? 
	`, userID).Find(&bots)
	return &bots
}
