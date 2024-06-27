package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/database/rdb"
	"github.com/wizzldev/chat/pkg/middlewares"
	"github.com/wizzldev/chat/pkg/repository"
)

type security struct{}

var Security security

func (security) Sessions(c *fiber.Ctx) error {
	sess, err := middlewares.Session(c)
	if err != nil {
		return err
	}
	sessID := sess.ID()

	sessions := repository.Session.AllForUser(authUserID(c))

	for _, s := range *sessions {
		if s.SessionID == sessID {
			s.Current = true
			break
		}
	}

	return c.JSON(sessions)
}

func (security) DestroySessions(c *fiber.Ctx) error {
	user := authUser(c)
	sessions := repository.Session.AllForUser(user.ID)

	var del []models.Session
	for _, session := range *sessions {
		err := rdb.Redis.Delete(session.SessionID)
		if err == nil {
			del = append(del, session)
		}
	}

	database.DB.Delete(&del)

	return c.JSON(fiber.Map{
		"status": "ok",
	})
}
