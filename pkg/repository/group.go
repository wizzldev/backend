package repository

import (
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/pkg/repository/paginator"
	"github.com/wizzldev/chat/pkg/utils/role"
	"strconv"
	"time"
)

type group struct{}

var Group group

func (group) FindPM(ID uint) *models.Group {
	return FindModelBy[models.Group]([]string{"id", "is_private_message"}, []any{ID, true})
}

func (group) Find(ID uint) *models.Group {
	return FindModelBy[models.Group]([]string{"id"}, []any{ID})
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
	select distinct users.id from groups
	inner join group_user on group_user.group_id = groups.id
	inner join users on users.id = group_user.user_id
	where groups.id = ?
    `, groupID).Find(&gIDs).Error

	if err != nil {
		log.Warn("Failed to execute query:", err)
	}

	return gIDs
}

func (group) IsGroupExists(userIDs [2]uint) (uint, bool) {
	var data struct {
		GroupID uint
	}

	database.DB.Raw(`
	select gu.group_id as group_id from group_user gu
	inner join users on gu.user_id = users.id
	inner join groups on groups.id = gu.group_id
	group by gu.group_id
	having sum(gu.user_id = ?) > 0
	and sum(gu.user_id = ?) > 0
	and count(*) = 2
	`, userIDs[0], userIDs[1]).
		Scan(&data)

	return data.GroupID, data.GroupID != 0
}

func (group) GetContactsForUser(userID uint, page int, authUser *models.User) *[]Contact {
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
		Verified         bool
		CustomInvite     *string
		UserID           uint
		SenderNickName   *string
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
       	groups.image_url,
		groups.verified,
		groups.custom_invite,
		groups.user_id,
		group_user.nick_name as sender_nick_name
	from messages
	join (
    	select group_id, max(created_at) as max_created_at from messages
    	group by group_id order by created_at desc
	) as latest_messages
	on messages.group_id = latest_messages.group_id
	and messages.created_at = latest_messages.max_created_at
	join users on messages.sender_id = users.id
	join groups on messages.group_id = groups.id
	join group_user on group_user.user_id = messages.sender_id and group_user.group_id = groups.id
	where groups.id in (
		select distinct group_user.group_id 
		from group_user 
		where group_user.user_id = ? and group_user.user_id
	)
	order by message_created_at desc limit 15 offset `+strconv.Itoa(offset)+`
	`, userID).Find(&data).Error

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
					break
				}
			}
			if groupName == "" {
				groupName = "You#allowTranslation"
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
			if imageURL == "" {
				imageURL = authUser.ImageURL
			}
		}

		contact := Contact{
			ID:               v.GroupID,
			Name:             groupName,
			ImageURL:         imageURL,
			Verified:         v.Verified,
			IsPrivateMessage: v.IsPrivateMessage,
			CustomInvite:     v.CustomInvite,
			CreatorID:        v.UserID,
			LastMessage: LastMessage{
				SenderID:   v.SenderID,
				SenderName: v.SenderName,
				NickName:   v.SenderNickName,
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
	_ = database.DB.Model(&models.Group{}).Preload("Theme").
		Where("id = ?", chatID).Find(&data).Error

	if data.IsPrivateMessage {
		var user models.User
		_ = database.DB.Raw(`
		select users.* from group_user
		inner join groups on groups.id = group_user.group_id
		inner join users on users.id = group_user.user_id
		where group_id = ? and users.id != ?
		limit 1
		`, data.ID, userID).Scan(&user).Error
		data.ImageURL = user.ImageURL
		if user.ID > 0 {
			data.Name = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
		}
	}

	return &data
}

func (group) GetUserRoles(gID uint, uID uint, roles role.Roles) role.Roles {
	var gUser models.GroupUser
	database.DB.Model(&models.GroupUser{}).Where("user_id = ? and group_id = ?", uID, gID).First(&gUser)

	for _, r := range gUser.Roles {
		realRole, err := role.New(r)
		if err != nil {
			continue
		}
		if realRole == role.Creator {
			return *role.All()
		}
		if realRole == role.Admin {
			roles = *role.All()
			roles.Revoke(role.Creator)
			break
		}
		roles = append(roles, realRole)
	}

	return roles
}

func (group) FindGroupUser(gID uint, uID uint) *models.GroupUser {
	var gUser models.GroupUser
	err := database.DB.Model(&models.GroupUser{}).Where("group_id = ? and user_id = ?", gID, uID).Find(&gUser)
	fmt.Println(err)
	return &gUser
}

func (group) IsBanned(groupID, userID uint) bool {
	var ban models.Ban
	database.DB.Model(&models.Ban{}).Where("user_id = ? and group_id = ?", userID, groupID).First(&ban)
	return ban.Exists()
}

func (group) IsGroupUserExists(groupID, userID uint) bool {
	var count int64
	database.DB.Model(&models.GroupUser{}).Where("group_id = ? and user_id = ?", groupID, userID).
		Limit(1).
		Count(&count)
	return count > 0
}

func (group) CustomInviteExists(s string) bool {
	var count int64
	database.DB.Model(&models.Group{}).Where("custom_invite not null and lower(custom_invite) = lower(?)", s).
		Limit(1).
		Count(&count)
	return count > 0
}

func (group) Users(gID uint, cursor string) (Pagination[models.User], error) {
	query := database.DB.Model(&models.User{}).
		Preload("GroupUser").
		Where("users.id in (select user_id from group_user where group_id = ?)", gID)

	data, next, prev, err := paginator.Paginate[models.User](query, &paginator.Config{
		Cursor:     cursor,
		Order:      "desc",
		Limit:      30,
		PointsNext: false,
	})

	return Pagination[models.User]{
		Data:       data,
		NextCursor: next,
		Previous:   prev,
	}, err
}

func (group) UserCount(gID uint) int {
	var count int64
	database.DB.Model(&models.User{}).
		Where("users.id in (select user_id from group_user where group_id = ?)", gID).
		Count(&count)

	return int(count)
}
