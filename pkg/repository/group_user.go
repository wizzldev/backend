package repository

import (
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
)

type groupUser struct{}

var GroupUser groupUser

func (groupUser) Delete(groupID uint, userID uint) {
	database.DB.Model(&models.GroupUser{}).Delete("group_id = ? and user_id = ?", groupID, userID)
}
