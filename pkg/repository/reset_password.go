package repository

import (
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
)

type resetPassword struct{}

var ResetPassword resetPassword

func (resetPassword) FindUserByToken(token string) *models.User {
	var user models.User
	database.DB.Raw(`
	select users.* from users
	inner join reset_passwords on users.id = reset_passwords.user_id
	where reset_passwords.token = ? 
	and reset_passwords.created_at > DATE_SUB(NOW(), INTERVAL 1 DAY)
	limit 1
	`, token).Scan(&user)
	return &user
}

func (resetPassword) FindLatestForUser(uid uint) *models.ResetPassword {
	var model models.ResetPassword
	database.DB.Model(&models.ResetPassword{}).Where("user_id = ? and created_at > DATE_SUB(NOW(), INTERVAL 1 DAY)", uid).Order("created_at DESC").First(&model)
	return &model
}
