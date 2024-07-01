package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"github.com/wizzldev/chat/app/events"
	"github.com/wizzldev/chat/app/requests"
	"github.com/wizzldev/chat/app/services"
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/pkg/repository"
	"github.com/wizzldev/chat/pkg/utils/role"
	"github.com/wizzldev/chat/pkg/ws"
)

type group struct {
	*services.Storage
	Cache *services.WSCache
}

var Group = &group{}

func (g *group) Init(store *services.Storage) {
	g.Storage = store
}

func (*group) New(c *fiber.Ctx) error {
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
		ImageURL:         "group.webp",
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

func (*group) GetInfo(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	g := repository.Group.GetChatUser(uint(id), authUserID(c))

	return c.JSON(g)
}

func (*group) GetAllRoles(c *fiber.Ctx) error {
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

func (g *group) UploadGroupImage(c *fiber.Ctx) error {
	serverID := c.Params("id")

	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	group := repository.Group.Find(uint(id))
	if group.ID < 1 || group.IsPrivateMessage {
		return fiber.ErrNotFound
	}

	fileH, err := c.FormFile("image")
	if err != nil {
		return err
	}

	file, err := g.Storage.StoreAvatar(fileH)
	if err != nil {
		return err
	}

	group.ImageURL = file.Discriminator + ".webp"
	database.DB.Save(group)

	err = events.DispatchMessage(serverID, g.Cache.GetGroupMemberIDs(serverID), uint(id), authUser(c), &ws.ClientMessage{
		Type:     "update.image",
		HookID:   c.Query("hook_id"),
		DataJSON: "{}",
	})

	if err != nil {
		return err
	}

	return c.JSON(group)
}
