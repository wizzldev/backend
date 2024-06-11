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
	where reset_passwords.token = ? limit 1
	`, token).Scan(&user)
	return &user
}
