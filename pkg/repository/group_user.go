package repository

import (
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
)

type groupUser struct{}

var GroupUser groupUser

func (groupUser) Delete(groupID uint, userID uint) {
	database.DB.Where("group_id = ? and user_id = ?", groupID, userID).Delete(&models.GroupUser{})
}

func (groupUser) Find(groupID uint, userID uint) (*models.GroupUser, error) {
	var gu models.GroupUser

	err := database.DB.Model(&models.GroupUser{}).Where("user_id = ? and group_id = ?", userID, groupID).First(&gu).Error

	return &gu, err
}

func (groupUser) Update(gu *models.GroupUser) error {
	return database.DB.Where("group_id = ? and user_id = ?", gu.GroupID, gu.UserID).First(&gu).Error
}
