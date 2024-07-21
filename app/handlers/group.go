package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"github.com/wizzldev/chat/app/events"
	"github.com/wizzldev/chat/app/requests"
	"github.com/wizzldev/chat/app/services"
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/pkg/configs"
	"github.com/wizzldev/chat/pkg/repository"
	"github.com/wizzldev/chat/pkg/utils/role"
	"github.com/wizzldev/chat/pkg/ws"
	"slices"
	"strings"
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
		ImageURL:         configs.DefaultGroupImage,
		Name:             data.Name,
		Roles:            roles,
		IsPrivateMessage: false,
		Users:            users,
	}

	database.DB.Create(&g)

	userID := authUserID(c)

	message := models.Message{
		HasGroup: models.HasGroup{
			GroupID: g.ID,
		},
		Type: "chat.create",
		HasMessageSender: models.HasMessageSender{
			SenderID: userID,
		},
	}
	database.DB.Create(&message)

	database.DB.Where("group_id = ? and user_id = ?", g.ID, userID).Save(&models.GroupUser{
		HasGroup: models.HasGroupID(g.ID),
		HasUser:  models.HasUserID(userID),
		Roles:    []string{string(role.Creator)},
	})

	return c.JSON(fiber.Map{
		"group_id": g.ID,
	})
}

func (*group) GetInfo(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	userID := authUserID(c)
	g := repository.Group.GetChatUser(uint(id), userID)

	return c.JSON(fiber.Map{
		"id":         g.ID,
		"created_at": g.CreatedAt,
		"updated_at": g.UpdatedAt,
		"image_url":  g.ImageURL,
		"name":       g.Name,
		"roles":      g.Roles,
		"your_roles": repository.Group.GetUserRoles(g.ID, userID, *role.NewRoles(g.Roles)),
	})
}

func (*group) GetAllRoles(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"roles": role.All(),
		"recommended": []role.Role{
			role.EditGroupImage,
			role.EditGroupName,
			role.EditGroupTheme,
			role.SendMessage,
			role.AttachFile,
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

	gr := repository.Group.Find(uint(id))
	if gr.ID < 1 || gr.IsPrivateMessage {
		return fiber.ErrBadRequest
	}

	fileH, err := c.FormFile("image")
	if err != nil {
		return err
	}

	file, err := g.Storage.StoreAvatar(fileH)
	if err != nil {
		return err
	}

	gr.ImageURL = file.Discriminator + ".webp"
	database.DB.Save(gr)

	if gr.ImageURL != configs.DefaultGroupImage {
		_ = g.Storage.RemoveByDisc(strings.SplitN(gr.ImageURL, ".", 2)[0])
	}

	err = events.DispatchMessage(serverID, g.Cache.GetGroupMemberIDs(serverID), uint(id), authUser(c), &ws.ClientMessage{
		Type:     "update.image",
		HookID:   c.Query("hook_id"),
		DataJSON: "{}",
	})

	if err != nil {
		return err
	}

	return c.JSON(gr)
}

func (g *group) ModifyRoles(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	serverID := c.Params("id")
	if err != nil {
		return err
	}

	gr := repository.Group.Find(uint(id))
	if gr.ID < 1 || gr.IsPrivateMessage {
		return fiber.ErrBadRequest
	}

	roles := validation[requests.ModifyRoles](c)

	userRoles := repository.Group.GetUserRoles(uint(id), authUserID(c), *role.NewRoles(gr.Roles))
	if !userRoles.Can(role.Creator) {
		if slices.Contains(gr.Roles, string(role.Creator)) != slices.Contains(roles.Roles, string(role.Creator)) {
			return fiber.ErrForbidden
		}
	}

	gr.Roles = roles.Roles

	database.DB.Save(gr)

	userIDs := g.Cache.GetGroupMemberIDs(serverID)

	_ = events.DispatchMessage(serverID, userIDs, uint(id), authUser(c), &ws.ClientMessage{
		Type:     "update.roles",
		DataJSON: "{}",
	})

	events.SendToGroup(serverID, userIDs, ws.Message{
		Event: "reload",
		Data:  nil,
	})

	return c.JSON(fiber.Map{
		"status": "ok",
	})
}

func (g *group) EditName(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	serverID := c.Params("id")
	if err != nil {
		return err
	}

	gr := repository.Group.Find(uint(id))
	if gr.ID < 1 || gr.IsPrivateMessage {
		return fiber.ErrBadRequest
	}

	data := validation[requests.EditGroupName](c)
	gr.Name = data.Name
	database.DB.Save(gr)

	userIDs := g.Cache.GetGroupMemberIDs(serverID)

	_ = events.DispatchMessage(serverID, userIDs, uint(id), authUser(c), &ws.ClientMessage{
		Type:     "update.name",
		DataJSON: "{}",
	})

	return c.JSON(fiber.Map{
		"status": "ok",
	})
}

func (*group) Delete(*fiber.Ctx) error {
	// TODO
	return nil
}
