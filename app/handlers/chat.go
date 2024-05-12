package handlers

import (
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/app/requests"
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/pkg/repository"
	"github.com/wizzldev/chat/pkg/ws"
	"strconv"
)

type chat struct{}

var Chat chat

func (chat) Contacts(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	data := repository.Group.GetContactsForUser(authUserID(c), page)
	return c.JSON(data)
}

func (chat) PrivateMessage(c *fiber.Ctx) error {
	requestedUserID, err := c.ParamsInt("id", 0)
	userID := uint(requestedUserID)

	user := authUser(c)
	if err != nil {
		return err
	}

	if repository.Block.IsBlocked(userID, user.ID) {
		return fiber.NewError(fiber.StatusForbidden, "You are blocked")
	}

	if gID, ok := repository.Group.IsGroupExists([]uint{user.ID, userID}); ok {
		return c.JSON(fiber.Map{
			"pm_id": gID,
		})
	}

	group := models.Group{
		Name:     nil,
		ImageURL: nil,
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
	database.DB.Create(&group)
	database.DB.Create(&models.Message{
		HasGroup: models.HasGroup{
			GroupID: group.ID,
		},
		HasMessageSender: models.HasMessageSender{
			SenderID: user.ID,
		},
		Type:     "chat.create",
		DataJSON: "{}",
	})

	return c.JSON(fiber.Map{
		"pm_id": group.ID,
	})
}

func (chat) Search(c *fiber.Ctx) error {
	v := validation[requests.SearchContacts](c)
	rawPage := c.Query("page", "1")
	page, err := strconv.Atoi(rawPage)

	if err != nil {
		return err
	}

	users := repository.User.Search(v.FirstName, v.LastName, v.Email, page)
	return c.JSON(users)
}

func (chat) Find(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id", 0)
	if err != nil {
		return err
	}

	g := repository.Group.GetChatUser(uint(id), authUserID(c))
	latestMessages := repository.Message.Latest(uint(id))
	return c.JSON(fiber.Map{
		"group":    g,
		"messages": latestMessages,
	})
}

func (chat) Connect(c *fiber.Ctx) error {
	serverID := c.Params("id")
	serverIDint, _ := strconv.Atoi(serverID)

	if !repository.Group.CanUserAccess(uint(serverIDint), authUser(c)) {
		return fmt.Errorf("you are not allowed to access this chat")
	}

	server, ok := ws.WebSocket[serverID]

	if !ok {
		server = ws.NewServer(serverID)
		ws.WebSocket[serverID] = server
	}

	return websocket.New(ws.WebSocket[serverID].AddConnection)(c)
}
