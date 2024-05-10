package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/app/requests"
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/pkg/middlewares"
	"github.com/wizzldev/chat/pkg/repository"
	"github.com/wizzldev/chat/pkg/utils"
)

type auth struct{}

var Auth auth

func (auth) Login(c *fiber.Ctx) error {
	sess, err := middlewares.Session(c)
	if err != nil {
		return err
	}
	defer sess.Save()

	loginRequest := validation[requests.Login](c)

	user := repository.User.FindByEmail(loginRequest.Email)
	if user.ID < 1 {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}

	if !utils.NewPassword(loginRequest.Password).Compare(user.Password) {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}

	sess.Set(utils.SessionAuthUserID, user.ID)
	return c.JSON(fiber.Map{
		"user":    user,
		"session": sess.ID(),
	})
}

func (auth) Register(c *fiber.Ctx) error {
	registerRequest := validation[requests.Register](c)

	if repository.User.IsEmailExists(registerRequest.Email) {
		return fiber.NewError(fiber.StatusBadRequest, "An account already exists with this email address")
	}

	password, err := utils.NewPassword(registerRequest.Password).Hash()

	if err != nil {
		return err
	}

	user := models.User{
		FirstName: registerRequest.FirstName,
		LastName:  registerRequest.LastName,
		Email:     registerRequest.Email,
		Password:  password,
	}

	err = database.DB.Create(&user).Error

	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Please verify your email before login",
	})
}

func (auth) Logout(c *fiber.Ctx) error {
	sess, err := middlewares.Session(c)

	if err != nil {
		return err
	}

	return sess.Destroy()
}
