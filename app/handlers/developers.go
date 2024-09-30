package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/app/requests"
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/pkg/repository"
	"github.com/wizzldev/chat/pkg/utils"
)

type developers struct{}

var Developers developers

func (developers) GetApplications(c *fiber.Ctx) error {
	bots := repository.Bot.FindBotsForUserID(authUserID(c))
	return c.JSON(fiber.Map{
		"bots": bots,
	})
}

func (developers) CreateApplication(c *fiber.Ctx) error {
	userID := authUserID(c)
	count := repository.Bot.CountForUser(userID)
	if count >= 3 {
		return fiber.NewError(fiber.StatusTooManyRequests, "A user cannot create more than 3 bots")
	}

	request := validation[requests.NewBot](c)

	token := utils.NewRandom().String(200)

	bot := models.User{
		FirstName: request.Name,
		Password:  token,
		ImageURL:  "default.webp",
		IsBot:     true,
	}

	err := database.DB.Create(&bot).Error
	if err != nil {
		return errors.New("Unknown error occurred when creating Bot")
	}

	bot.EnableIPCheck = false
	database.DB.Save(&bot)

	botUser := &models.UserBot{
		HasUser: models.HasUserID(userID),
		HasBot:  models.HasBotID(bot.ID),
	}
	database.DB.Create(&botUser)

	return c.JSON(fiber.Map{
		"token": token,
	})
}

func (developers) RegenerateApplicationToken(c *fiber.Ctx) error {
	rawID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	userID := authUserID(c)
	botID := uint(rawID)

	bot := repository.Bot.FindUserBot(userID, botID)
	if !bot.Exists() {
		return fiber.ErrNotFound
	}

	token := utils.NewRandom().String(200)
	bot.Password = token

	return c.JSON(fiber.Map{
		"token": token,
	})
}
