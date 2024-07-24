package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/app/events"
	"github.com/wizzldev/chat/app/services"
	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/pkg/repository"
	"github.com/wizzldev/chat/pkg/utils/role"
	"github.com/wizzldev/chat/pkg/ws"
	"strconv"
)

type groupHelpers struct{}

func (groupHelpers) GetAllRoles(c *fiber.Ctx) error {
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

func (groupHelpers) sendMessage(cache *services.WSCache, gID uint, user *models.User, message *ws.ClientMessage) []uint {
	serverID := strconv.Itoa(int(gID))
	userIDs := cache.GetGroupMemberIDs(serverID)
	_ = events.DispatchMessage(serverID, userIDs, gID, user, message)
	return userIDs
}

func (groupHelpers) group(id string) (*models.Group, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	g := repository.Group.Find(uint(idInt))
	if g.ID < 1 {
		return nil, errors.New("group does not exits")
	}
	return g, nil
}
