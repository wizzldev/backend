package repository

import (
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
)

type group struct{}

var Group group

func (group) FindPM(ID uint) *models.Group {
	return FindModelBy[models.Group]([]string{"id", "is_private_message"}, []any{ID, true})
}

func (group) CanUserAccess(groupID uint, u *models.User) bool {
	var count int64
	err := database.DB.Raw(`
	select count(*) from groups
	inner join group_user on group_user.group_id = groups.id
	inner join users on users.id = group_user.user_id
	where groups.id = ? and users.id = ?
	limit 1
	`, groupID, u.ID).
		Count(&count).Error

	if err != nil {
		log.Warn("Failed to execute query:", err)
		return false
	}

	return count > 0
}

func (group) GetUserIDs(groupID uint) []uint {
	var gIDs []uint
	err := database.DB.Raw(`
	select users.id from groups
	inner join group_user on group_user.group_id = groups.id
	inner join users on users.id = group_user.user_id
	where groups.id = ?
    `, groupID).Find(&gIDs).Error

	if err != nil {
		log.Warn("Failed to execute query:", err)
	}

	return gIDs
}

func (group) IsGroupExists(userIDs []uint) (uint, bool) {
	var data struct {
		MemberCount int64
		GroupID     uint
	}

	err := database.DB.Raw(`
	select count(distinct users.id) member_count, groups.id group_id from groups
	inner join group_user on group_user.group_id = groups.id
	inner join users on users.id = group_user.user_id
	where users.id in (?)
	order by groups.updated_at
	limit 1
	`, userIDs).
		Scan(&data).Error

	if err != nil {
		log.Warn("Failed to execute query:", err)
		return 0, true
	}

	fmt.Println(data)

	return data.GroupID, int(data.MemberCount) == len(userIDs)
}
