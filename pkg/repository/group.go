package repository

import (
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
	"strconv"
	"time"
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

func (group) GetContactsForUser(userID uint, page int) *[]Contact {
	var perPage = 15
	var offset = perPage * (page - 1)

	var data []struct {
		MessageContent   *string
		MessageType      string
		MessageCreatedAt time.Time
		SenderID         uint
		SenderName       string
		GroupID          uint
		IsPrivateMessage bool
		GroupName        *string
		ImageURL         *string
	}
	_ = database.DB.Raw(`
	select 
		messages.content as message_content,
		messages.type as message_type,
		messages.created_at as message_created_at,
       	users.id as sender_id,
       	users.first_name as sender_name,
       	groups.id as group_id,
       	groups.is_private_message,
       	groups.name as group_name,
       	groups.image_url
	from messages
	join (
    	select group_id, max(created_at) as max_created_at from messages
    	group by group_id order by created_at desc
	) as latest_messages
	on messages.group_id = latest_messages.group_id
	and messages.created_at = latest_messages.max_created_at
	join users on messages.sender_id = users.id
	join groups on messages.group_id = groups.id
	where groups.id in (
		select distinct group_user.group_id from group_user where group_user.user_id = ?
	)
	order by message_created_at desc limit 15 offset `+strconv.Itoa(offset)+`
	`, userID).Find(&data).Error

	fmt.Println(data)

	var privateMessageIDs []uint
	for _, v := range data {
		if v.IsPrivateMessage {
			privateMessageIDs = append(privateMessageIDs, v.GroupID)
		}
	}

	var userGroupMap []struct {
		GroupID       uint
		UserFirstName string
		UserLastName  string
		UserImageUrl  string
	}
	_ = database.DB.Raw(`
	select 
		groups.id as group_id,
		users.first_name as user_first_name, 
		users.last_name as user_last_name,
		users.image_url as user_image_url
	from group_user
	inner join groups on groups.id = group_user.group_id
	inner join users on users.id = group_user.user_id
	where group_user.group_id in (?) and users.id != ?
	`, privateMessageIDs, userID).Find(&userGroupMap).Error

	var contacts []Contact

	for _, v := range data {
		groupName := ""
		imageURL := ""

		if v.GroupName != nil {
			groupName = *v.GroupName
		} else {
			for _, u := range userGroupMap {
				if u.GroupID == v.GroupID {
					groupName = fmt.Sprintf("%s %s", u.UserFirstName, u.UserLastName)
				}
			}
		}

		if v.ImageURL != nil {
			imageURL = *v.ImageURL
		} else {
			for _, u := range userGroupMap {
				if u.GroupID == v.GroupID {
					imageURL = u.UserImageUrl
				}
			}
		}

		contact := Contact{
			ID:       v.GroupID,
			Name:     groupName,
			ImageURL: imageURL,
			LastMessage: LastMessage{
				SenderID:   v.SenderID,
				SenderName: v.SenderName,
				Content:    v.MessageContent,
				Type:       v.MessageType,
				Date:       v.MessageCreatedAt,
			},
		}
		contacts = append(contacts, contact)
	}

	return &contacts
}

func (group) GetChatUser(chatID uint, userID uint) *models.Group {
	var data models.Group
	_ = database.DB.Model(&models.Group{}).
		Where("id = ?", chatID).Find(&data).Error

	if data.IsPrivateMessage {
		var user models.User
		_ = database.DB.Raw(`
		select users.* from group_user
		inner join groups on groups.id = group_user.group_id
		inner join users on users.id = group_user.user_id
		where group_id = ? and users.id != ? limit 1
		`, data.ID, userID).Scan(&user).Error
		data.ImageURL = &user.ImageURL
		data.Name = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	}

	return &data
}
