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
	"strconv"
)

type invite struct {
	cache *services.WSCache
}

var Invite invite

func (invite) Create(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	data := validation[requests.NewInvite](c)

	i := &models.Invite{
		HasUser:    models.HasUserID(authUserID(c)),
		HasGroup:   models.HasGroupID(uint(id)),
		Key:        repository.Invite.CreateCode(),
		Expiration: data.Expiration,
	}

	if data.MaxUsage > 0 {
		i.MaxUsage = &data.MaxUsage
	}

	database.DB.Create(i)

	return c.JSON(i)
}

func (i invite) Use(c *fiber.Ctx) error {
	userID := authUserID(c)

	groupID, err := i.getGroupID(c.Params("code", ""), userID)
	if err != nil {
		return err
	}

	gu := &models.GroupUser{
		HasGroup: models.HasGroupID(groupID),
		HasUser:  models.HasUserID(userID),
	}
	database.DB.Save(gu)

	serverID := strconv.Itoa(int(groupID))
	_ = events.DispatchUserJoin(serverID, i.cache.GetGroupMemberIDs(serverID), authUser(c), groupID)

	return c.JSON(fiber.Map{
		"status": "success",
	})
}

func (invite) getGroupID(code string, userID uint) (uint, error) {
	inv := repository.Invite.FindInviteByCode(code)
	fmt.Println("code", code, inv)
	if inv.IsValid() {
		if repository.Group.IsBanned(inv.GroupID, userID) {
			return 0, fiber.ErrForbidden
		}

		inv.Decrement()
		database.DB.Save(inv)
		return inv.GroupID, nil
	}

	g := repository.Invite.FindGroupInviteByCode(code)
	if !g.Exists() {
		return 0, fiber.ErrNotFound
	}

	if repository.Group.IsBanned(g.ID, userID) {
		return 0, fiber.ErrForbidden
	}

	return g.ID, nil
}
