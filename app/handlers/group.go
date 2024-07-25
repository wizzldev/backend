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
	"strconv"
	"strings"
)

type group struct {
	groupHelpers
	*services.Storage
	Cache *services.WSCache
}

var Group = &group{}

func (g *group) Init(store *services.Storage, cache *services.WSCache) {
	g.Storage = store
	g.Cache = cache
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

	userID := authUserID(c)

	g := models.Group{
		ImageURL:         configs.DefaultGroupImage,
		Name:             data.Name,
		Roles:            roles,
		IsPrivateMessage: false,
		Users:            users,
		HasUser:          models.HasUserID(userID),
	}

	database.DB.Create(&g)

	message := models.Message{
		HasGroup:         models.HasGroupID(g.ID),
		Type:             "chat.create",
		HasMessageSender: models.HasMessageSenderID(userID),
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
		"id":                 g.ID,
		"created_at":         g.CreatedAt,
		"updated_at":         g.UpdatedAt,
		"image_url":          g.ImageURL,
		"name":               g.Name,
		"roles":              g.Roles,
		"is_private_message": g.IsPrivateMessage,
		"is_verified":        g.Verified,
		"custom_invite":      g.CustomInvite,
		"emoji":              g.Emoji,
		"your_roles":         repository.Group.GetUserRoles(g.ID, userID, *role.NewRoles(g.Roles)),
	})
}

func (g *group) UploadGroupImage(c *fiber.Ctx) error {
	gr, err := g.group(c.Params("id"))
	if err != nil {
		return err
	}

	if gr.IsPrivateMessage {
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

	// FIX: firstly delete the image then save the new
	if gr.ImageURL != configs.DefaultGroupImage {
		_ = g.Storage.RemoveByDisc(strings.SplitN(gr.ImageURL, ".", 2)[0])
	}

	gr.ImageURL = file.Discriminator + ".webp"
	database.DB.Save(gr)

	g.sendMessage(g.Cache, gr.ID, authUser(c), &ws.ClientMessage{
		Type:     "update.image",
		HookID:   c.Query("hook_id"),
		DataJSON: "{}",
	})

	return c.JSON(gr)
}

func (g *group) ModifyRoles(c *fiber.Ctx) error {
	serverID := c.Params("id")
	gr, err := g.group(serverID)
	if err != nil {
		return err
	}

	if gr.IsPrivateMessage {
		return fiber.ErrBadRequest
	}

	roles := validation[requests.ModifyRoles](c)

	userRoles := repository.Group.GetUserRoles(gr.ID, authUserID(c), *role.NewRoles(gr.Roles))
	if !userRoles.Can(role.Creator) {
		if slices.Contains(gr.Roles, string(role.Creator)) != slices.Contains(roles.Roles, string(role.Creator)) {
			return fiber.ErrForbidden
		}
	}

	gr.Roles = roles.Roles

	database.DB.Save(gr)

	userIDs := g.sendMessage(g.Cache, gr.ID, authUser(c), &ws.ClientMessage{
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
	gr, err := g.group(c.Params("id"))
	if err != nil {
		return err
	}

	if gr.IsPrivateMessage {
		return fiber.ErrBadRequest
	}

	data := validation[requests.EditGroupName](c)
	gr.Name = data.Name
	database.DB.Save(gr)

	g.sendMessage(g.Cache, gr.ID, authUser(c), &ws.ClientMessage{
		Type:     "update.name",
		DataJSON: "{}",
	})

	return c.JSON(fiber.Map{
		"status": "ok",
	})
}

func (g *group) CustomInvite(c *fiber.Ctx) error {
	gr, err := g.group(c.Params("id"))
	if err != nil {
		return err
	}

	data := validation[requests.CustomInvite](c)

	if gr.IsPrivateMessage {
		return fiber.ErrBadRequest
	}

	if repository.Group.CustomInviteExists(data.Invite) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status": "already-exists",
		})
	}

	gr.CustomInvite = &data.Invite
	database.DB.Save(gr)

	return c.JSON(fiber.Map{
		"status": "ok",
	})
}

func (g *group) Leave(c *fiber.Ctx) error {
	gr, err := g.group(c.Params("id"))
	if err != nil {
		return err
	}

	userID := authUserID(c)

	if gr.IsPrivateMessage || gr.UserID == userID {
		return fiber.ErrBadRequest
	}

	repository.GroupUser.Delete(gr.ID, userID)

	g.sendMessage(g.Cache, gr.ID, authUser(c), &ws.ClientMessage{
		Type: "leave",
	})

	return c.JSON(fiber.Map{
		"status": "ok",
	})
}

func (g *group) Delete(c *fiber.Ctx) error {
	gr, err := g.group(c.Params("id"))
	if err != nil {
		return err
	}

	if gr.ImageURL != configs.DefaultGroupImage {
		_ = g.Storage.RemoveByDisc(strings.SplitN(gr.ImageURL, ".", 2)[0])
	}

	g.sendMessage(g.Cache, gr.ID, authUser(c), &ws.ClientMessage{
		Type:    "delete",
		Content: strconv.Itoa(int(gr.ID)),
	})

	worker := &models.Worker{
		Command: "cleanup.group",
		Data:    strconv.Itoa(int(gr.ID)),
	}
	database.DB.Create(&worker)

	return c.JSON(fiber.Map{
		"status": "ok",
	})
}

func (g *group) Emoji(c *fiber.Ctx) error {
	gr, err := g.group(c.Params("id"))
	if err != nil {
		return err
	}

	data := validation[requests.Emoji](c)
	gr.Emoji = &data.Emoji
	database.DB.Save(gr)

	return c.JSON(fiber.Map{
		"status": "ok",
	})
}

func (g *group) Users(c *fiber.Ctx) error {
	gID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	data, err := repository.Group.Users(uint(gID), c.Query("cursor"))
	if err != nil {
		return err
	}

	return c.JSON(data)
}

func (g *group) UserCount(c *fiber.Ctx) error {
	gID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"count": repository.Group.UserCount(uint(gID)),
	})
}
