package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/database/models"
	"github.com/wizzldev/chat/pkg/configs"
	"github.com/wizzldev/chat/pkg/middlewares"
	"github.com/wizzldev/chat/pkg/repository"
	"github.com/wizzldev/chat/pkg/utils"
	"net"
)

type Auth struct {
	ctx *fiber.Ctx
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"-"`
}

type AuthResponse struct {
	MustVerifyIP bool
	Token        string
	User         *models.User
}

func NewAuth(c *fiber.Ctx) *Auth {
	return &Auth{ctx: c}
}

func (a *Auth) Login(request *AuthRequest) (*AuthResponse, error) {
	sess, err := middlewares.Session(a.ctx)
	if err != nil {
		return nil, err
	}

	user := repository.User.FindByEmail(request.Email)
	if user.ID < 1 {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}

	if user.IsBot {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Cannot use bot as user")
	}

	if !utils.NewPassword(request.Password).Compare(user.Password) {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
	}

	if user.EmailVerifiedAt == nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Please verify your email before login")
	}

	ip := a.ctx.IP()
	if user.EnableIPCheck && !repository.User.IsIPAllowed(user.ID, ip) && !net.ParseIP(ip).IsPrivate() {
		return &AuthResponse{
			MustVerifyIP: true,
			User:         user,
		}, nil
	}

	sess.Set(configs.SessionAuthUserID, user.ID)
	sessID := sess.ID()
	err = sess.Save()
	if err != nil {
		return nil, err
	}

	database.DB.Create(&models.Session{
		HasUser:   models.HasUserID(user.ID),
		IP:        ip,
		SessionID: sessID,
		Agent:     string(a.ctx.Request().Header.Peek("User-Agent")),
	})

	return &AuthResponse{
		MustVerifyIP: false,
		Token:        sessID,
		User:         user,
	}, nil
}
