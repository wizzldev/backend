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
