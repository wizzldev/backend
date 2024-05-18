package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"github.com/wizzldev/chat/app/requests"
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/pkg/repository"
	"github.com/wizzldev/chat/pkg/utils/role"
)

type group struct{}

var Group group

func (group) New(c *fiber.Ctx) error {
	data := validation[requests.NewGroup](c)

	userIDs := repository.IDsExists[models.User](data.UserIDs)
	var users []*models.User

	users = append(users, &models.User{
		Base: models.Base{ID: authUserID(c)},
	})
	for _, id := range userIDs {
		users = append(users, &models.User{Base: models.Base{ID: id}})
	}

	var roles pq.StringArray
	for _, r := range data.Roles {
		roles = append(roles, r)
	}

	g := models.Group{
		ImageURL:         "https://images.unsplash.com/photo-1493612276216-ee3925520721?q=80&w=1000&auto=format&fit=crop&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8Mnx8cmFuZG9tfGVufDB8fDB8fHww",
		Name:             data.Name,
		Roles:            roles,
		IsPrivateMessage: false,
		Users:            users,
	}

	database.DB.Create(&g)
	message := models.Message{
		HasGroup: models.HasGroup{
			GroupID: g.ID,
		},
		Type: "chat.create",
		HasMessageSender: models.HasMessageSender{
			SenderID: authUserID(c),
		},
	}
	database.DB.Create(&message)

	return c.JSON(fiber.Map{
		"group_id": g.ID,
	})
}

func (group) GetInfo(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	g := repository.Group.GetChatUser(uint(id), authUserID(c))

	return c.JSON(g)
}

func (group) GetAllRoles(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"roles": role.All(),
		"recommended": []role.Role{
			role.EditGroupImage,
			role.EditGroupName,
			role.EditGroupTheme,
			role.SendMessage,
			role.DeleteMessage,
			role.CreateIntegration,
			role.KickUser,
			role.InviteUser,
		},
	})
}
