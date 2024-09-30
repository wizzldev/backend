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

func (bot) FindByID(id uint) *models.User {
	return FindModelBy[models.User]([]string{"id", "is_bot"}, []any{id, true})
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

func (bot) CountForUser(userID uint) int {
	var count int64
	database.DB.Model(&models.UserBot{}).Where("user_id = ?", userID).Count(&count)
	return int(count)
}

func (bot) FindUserBot(userID, botID uint) *models.User {
	var bot models.User
	database.DB.Raw(`
		select users.* from users
		inner join user_bots on user_bots.user_id = ? and user_bots.bot_id = ?
		limit 1
	`,
		userID, botID).First(&bot)
	return &bot
}
