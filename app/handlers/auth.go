package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/app/requests"
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/pkg/configs"
	"github.com/wizzldev/chat/pkg/middlewares"
	"github.com/wizzldev/chat/pkg/repository"
	"github.com/wizzldev/chat/pkg/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
		ImageURL:  "https://xsgames.co/randomusers/assets/avatars/female/71.jpg",
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

func (auth) RequestNewPassword(c *fiber.Ctx) error {
	newPasswordRequest := validation[requests.NewPassword](c)

	user := repository.User.FindByEmail(newPasswordRequest.Email)

	if user.ID < 1 {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	token := utils.NewRandom().String(30)
	resetPassword := models.ResetPassword{
		HasUser: models.HasUser{
			UserID: user.ID,
		},
		Token: token,
	}

	err := database.DB.Create(&resetPassword).Error

	if err != nil {
		return err
	}

	go func() {
		resetPasswordURL := fmt.Sprintf("%s/reset-password/%s", configs.Env.Frontend, token)
		mail := utils.NewMail(configs.Env.Email.SenderAddress, user.Email)
		mail.Subject("Reset your password")
		mail.TemplateBody("reset-password", map[string]string{
			"firstName":        cases.Title(language.English).String(user.FirstName),
			"resetPasswordURL": resetPasswordURL,
		}, fmt.Sprintf("Click <a href=\"%s\">here</a> to reset your password", resetPasswordURL))
		err := mail.Send()
		fmt.Println("Email sent with err:", err)
	}()

	return c.JSON(fiber.Map{
		"message": "New password has been sent to your email",
	})
}

func (auth) SetNewPassword(c *fiber.Ctx) error {
	newPasswordRequest := validation[requests.SetNewPassword](c)
	token := c.Params("token")
	user := repository.ResetPassword.FindUserByToken(token)

	pass, err := utils.NewPassword(newPasswordRequest.Password).Hash()

	if err != nil {
		return err
	}

	user.Password = pass
	database.DB.Save(&user)

	database.DB.Delete(&models.ResetPassword{Token: token})

	return c.JSON(fiber.Map{
		"message": "Password has been updated successfully",
	})
}
