package repository

import (
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
)

type emailVerification struct{}

var EmailVerification emailVerification

func (emailVerification) FindUserByToken(token string) *models.User {
	var user models.User
	database.DB.Raw(`
	select users.* from users
	inner join email_verifications on users.id = email_verifications.user_id
	where email_verifications.token = ? 
	and email_verifications.created_at > DATE_SUB(NOW(), INTERVAL 1 DAY)
	limit 1
	`, token).Scan(&user)
	return &user
}

func (emailVerification) FindLatestForUser(uid uint) *models.EmailVerification {
	var model models.EmailVerification
	database.DB.Model(&models.EmailVerification{}).Where("user_id = ? and created_at > DATE_SUB(NOW(), INTERVAL 1 DAY)", uid).Order("created_at DESC").First(&model)
	return &model
}
