package handlers

import (
	"fmt"
	"github.com/wizzldev/chat/pkg/services"
	"strconv"
	"time"

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

func (a auth) Login(c *fiber.Ctx) error {
	loginRequest := validation[requests.Login](c)

	service := services.NewAuth(c)
	data, err := service.Login(&services.AuthRequest{
		Email:    loginRequest.Email,
		Password: loginRequest.Password,
	})

	if err != nil {
		return err
	}

	if data.MustVerifyIP {
		a.sendIPVerification(data.User, c.IP())
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"show_ip_modal": true,
		})
	}

	return c.JSON(fiber.Map{
		"user":    data.User,
		"session": "Bearer " + data.Token,
	})
}

func (a auth) Register(c *fiber.Ctx) error {
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
		ImageURL:  configs.DefaultUserImage,
	}

	err = database.DB.Create(&user).Error
	if err != nil {
		return err
	}

	err = a.sendVerificationEmail(&user)
	if err != nil {
		return err
	}

	database.DB.Create(&models.AllowedIP{
		HasUser: models.HasUserID(user.ID),
		IP:      c.IP(),
		Active:  true,
	})

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"show_verification_modal": true,
	})
}

func (auth) Logout(c *fiber.Ctx) error {
	sess, err := middlewares.Session(c)

	if err != nil {
		return err
	}

	database.DB.Where("session_id = ?", sess.ID()).Delete(&models.Session{})

	return sess.Destroy()
}

func (a auth) RequestNewEmailVerification(c *fiber.Ctx) error {
	email := validation[requests.Email](c)
	user := repository.User.FindByEmail(email.Email)
	if user.ID < 1 {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	if user.EmailVerifiedAt != nil {
		return fiber.NewError(fiber.StatusConflict, "Your email address is already verified")
	}

	emailVerification := repository.EmailVerification.FindLatestForUser(user.ID)
	if emailVerification.ID > 0 {
		return fiber.NewError(fiber.StatusConflict, "Email verification request already sent")
	}

	err := a.sendVerificationEmail(user)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"status": "success",
	})
}

func (auth) VerifyEmail(c *fiber.Ctx) error {
	token := c.Params("token")
	user := repository.EmailVerification.FindUserByToken(token)

	if user.ID < 1 {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	now := time.Now()
	user.EmailVerifiedAt = &now

	err := database.DB.Save(&user).Error
	if err != nil {
		return err
	}

	database.DB.Where("token = ?", token).Delete(&models.EmailVerification{})

	return c.JSON(fiber.Map{
		"message": "Password has been updated successfully",
	})
}

func (a auth) RequestNewPassword(c *fiber.Ctx) error {
	newPasswordRequest := validation[requests.NewPassword](c)

	user := repository.User.FindByEmail(newPasswordRequest.Email)
	if user.ID < 1 {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	resetPassword := repository.ResetPassword.FindLatestForUser(user.ID)
	if resetPassword.ID > 0 {
		return fiber.NewError(fiber.StatusConflict, "Reset password request already sent")
	}

	err := a.sendResetPasswordEmail(user)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "New password has been sent to your email",
	})
}

func (auth) SetNewPassword(c *fiber.Ctx) error {
	newPasswordRequest := validation[requests.SetNewPassword](c)
	token := c.Params("token")

	user := repository.ResetPassword.FindUserByToken(token)
	if user.ID < 1 {
		return fiber.NewError(fiber.StatusNotFound, "Invalid or expired token")
	}

	pass, err := utils.NewPassword(newPasswordRequest.Password).Hash()

	if err != nil {
		return err
	}

	user.Password = pass
	database.DB.Save(&user)

	database.DB.Where("token = ?", token).Delete(&models.ResetPassword{})

	return c.JSON(fiber.Map{
		"message": "Password has been updated successfully",
	})
}

func (auth) IsPasswordResetExists(c *fiber.Ctx) error {
	token := c.Params("token")
	user := repository.ResetPassword.FindUserByToken(token)

	if user.ID < 1 {
		return fiber.NewError(fiber.StatusNotFound, "Reset password not found")
	}

	return c.JSON(fiber.Map{
		"exists": true,
	})
}

func (auth) AllowIP(c *fiber.Ctx) error {
	token := c.Params("token")
	var ip models.AllowedIP
	database.DB.Model(&models.AllowedIP{}).Where("verification = ?", token).First(&ip)
	if ip.ID < 1 {
		return fiber.NewError(fiber.StatusNotFound, "IP not found")
	}
	ip.Verification = ""
	ip.Active = true
	database.DB.Save(&ip)
	return c.JSON(fiber.Map{
		"allowed": true,
	})
}

// helpers
func (auth) sendVerificationEmail(user *models.User) error {
	token := strconv.Itoa(int(user.ID)) + utils.NewRandom().String(30)

	err := database.DB.Create(&models.EmailVerification{
		HasUser: models.HasUserID(user.ID),
		Token:   token,
	}).Error

	if err != nil {
		return err
	}

	go func() {
		resetPasswordURL := fmt.Sprintf("%s/verify-email/%s", configs.Env.Frontend, token)
		mail := utils.NewMail(configs.Env.Email.SenderAddress, user.Email)
		mail.Subject("Verify your email address")
		mail.TemplateBody("register", map[string]string{
			"firstName":       cases.Title(language.English).String(user.FirstName),
			"verificationURL": resetPasswordURL,
		}, fmt.Sprintf("Click <a href=\"%s\">here</a> to verify your email address", resetPasswordURL))
		err := mail.Send()
		fmt.Println("Email sent with err:", err)
	}()

	return nil
}

func (auth) sendResetPasswordEmail(user *models.User) error {
	token := strconv.Itoa(int(user.ID)) + utils.NewRandom().String(30)
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

	return nil
}

func (auth) sendIPVerification(user *models.User, ip string) {
	token := strconv.Itoa(int(user.ID)) + utils.NewRandom().String(30)
	ipVerification := models.AllowedIP{
		HasUser:      models.HasUserID(user.ID),
		IP:           ip,
		Active:       false,
		Verification: token,
	}
	database.DB.Create(&ipVerification)

	go func() {
		verifyIPURL := fmt.Sprintf("%s/ip-verification/%s", configs.Env.Frontend, token)
		mail := utils.NewMail(configs.Env.Email.SenderAddress, user.Email)
		mail.Subject("IP Verification")
		mail.TemplateBody("ip", map[string]string{
			"firstName":   cases.Title(language.English).String(user.FirstName),
			"ip":          ip,
			"verifyIPURL": verifyIPURL,
		}, fmt.Sprintf("Click <a href=\"%s\">here</a> to verify that IP address", verifyIPURL))
		err := mail.Send()
		fmt.Println("Email sent with err:", err)
	}()
}
