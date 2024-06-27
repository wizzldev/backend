package repository

import (
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
)

type ips struct{}

var IPs ips

func (ips) AllForUser(uID uint) []*models.AllowedIP {
	var ip []*models.AllowedIP
	database.DB.Model(&models.AllowedIP{}).Where("user_id = ? and active = ?", uID, true).Order("created_at desc").Limit(30).Find(&ip)
	return ip
}

func (ips) FindForUser(uID uint, id uint) *models.AllowedIP {
	var s *models.AllowedIP
	database.DB.Model(&models.AllowedIP{}).Where("user_id = ? and id = ? and active = ?", uID, id, true).First(&s)
	return s
}
