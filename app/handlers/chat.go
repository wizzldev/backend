package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/app/events"
	"github.com/wizzldev/chat/app/requests"
	"github.com/wizzldev/chat/app/services"
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/pkg/repository"
	"github.com/wizzldev/chat/pkg/utils"
	"github.com/wizzldev/chat/pkg/utils/role"
	"github.com/wizzldev/chat/pkg/ws"
	"net/url"
	"strconv"
)

type chat struct {
	*services.Storage
	Cache *services.WSCache
}

var Chat = &chat{}

func (ch *chat) Init(store *services.Storage, wsCache *services.WSCache) {
	ch.Storage = store
	ch.Cache = wsCache
}

func (*chat) Contacts(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	data := repository.Group.GetContactsForUser(authUserID(c), page, authUser(c))
	return c.JSON(data)
}

func (*chat) PrivateMessage(c *fiber.Ctx) error {
	requestedUserID, err := c.ParamsInt("id", 0)
	if err != nil {
		return err
	}

	userID := uint(requestedUserID)
	user := authUser(c)

	if repository.Block.IsBlocked(userID, user.ID) {
		return fiber.NewError(fiber.StatusForbidden, "You are blocked")
	}

	if gID, ok := repository.Group.IsGroupExists([2]uint{user.ID, userID}); ok {
		return c.JSON(fiber.Map{
			"pm_id": gID,
		})
	}

	g := models.Group{
		Users: []*models.User{
			{
				Base: models.Base{
					ID: userID,
				},
			},
			user,
		},
		IsPrivateMessage: true,
	}

	database.DB.Create(&g)
	database.DB.Create(&models.Message{
		HasGroup: models.HasGroup{
			GroupID: g.ID,
		},
		HasMessageSender: models.HasMessageSender{
			SenderID: user.ID,
		},
		Type:     "chat.create",
		DataJSON: "{}",
	})

	return c.JSON(fiber.Map{
		"pm_id": g.ID,
	})
}

func (*chat) Search(c *fiber.Ctx) error {
	v := validation[requests.SearchContacts](c)
	rawPage := c.Query("page", "1")
	page, err := strconv.Atoi(rawPage)

	if err != nil {
		return err
	}

	users := repository.User.Search(v.FirstName, v.LastName, v.Email, page)
	return c.JSON(users)
}

func (*chat) Find(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id", 0)
	if err != nil {
		return err
	}

	user := authUser(c)

	var isYourProfile = false

	g := repository.Group.GetChatUser(uint(id), authUserID(c))
	if g.ImageURL == "" && g.Name == "" {
		g.ImageURL = user.ImageURL
		g.Name = "You#allowTranslation"
		isYourProfile = true
	}

	roles := role.Roles{}
	roles = append(roles, repository.Group.GetUserRoles(g.ID, authUserID(c), *role.NewRoles(g.Roles))...)

	pagination, err := repository.Message.CursorPaginate(uint(id), c.Query("cursor"))
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"group":           g,
		"messages":        pagination,
		"user_roles":      roles,
		"is_your_profile": isYourProfile,
	})
}

func (*chat) Messages(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id", 0)
	if err != nil {
		return err
	}

	pagination, err := repository.Message.CursorPaginate(uint(id), c.Query("cursor"))
	if err != nil {
		return err
	}

	return c.JSON(pagination)
}

func (*chat) FindMessage(c *fiber.Ctx) error {
	id, err := c.ParamsInt("messageID")
	if err != nil {
		return err
	}

	return c.JSON(repository.Message.FindOne(uint(id)))
}

func (ch *chat) UploadFile(c *fiber.Ctx) error {
	serverID := c.Params("id")
	gID, err := strconv.Atoi(serverID)
	if err != nil {
		return err
	}

	fileH, err := c.FormFile("file")
	if err != nil {
		return err
	}

	token := utils.NewRandom().String(50)
	file, err := ch.Store(fileH, token)
	if err != nil {
		return err
	}

	user := authUser(c)

	err = events.DispatchMessage(serverID, ch.Cache.GetGroupMemberIDs(serverID), uint(gID), user, &ws.ClientMessage{
		Content:  "none",
		Type:     "file:" + file.Type,
		DataJSON: fmt.Sprintf(`{"fetchFrom": "/storage/files/%s/%s", "hasAccessToken": true, "accessToken": "%s"}`, file.Discriminator, url.QueryEscape(file.Name), token),
		HookID:   c.Query("hook_id"),
	})
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"status": "success",
	})
}
