package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/app/handlers"
	"github.com/wizzldev/chat/app/requests"
	"github.com/wizzldev/chat/pkg/middlewares"
	"github.com/wizzldev/chat/pkg/utils/role"
	"time"
)

func RegisterAPI(r fiber.Router) {
	r.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./pkg/configs/build.json")
	})

	{
		r.Post("/login", requests.Use[requests.Login](), middlewares.NewAuthLimiter(), handlers.Auth.Login)
		r.Get("/allow-ip/:token", handlers.Auth.AllowIP)
		r.Post("/register", requests.Use[requests.Register](), handlers.Auth.Register)
		r.Get("/logout", handlers.Auth.Logout)
		r.Post("/request-new-password", requests.Use[requests.NewPassword](), handlers.Auth.RequestNewPassword)
		r.Get("/set-new-password/:token", handlers.Auth.IsPasswordResetExists)
		r.Post("/set-new-password/:token", requests.Use[requests.SetNewPassword](), handlers.Auth.SetNewPassword)
		r.Get("/verify-email/:token", handlers.Auth.VerifyEmail)
		r.Post("/request-new-email-verification", requests.Use[requests.Email](), handlers.Auth.RequestNewEmailVerification)
	}

	auth := r.Group("/", middlewares.AnyAuth)
	{
		auth.Get("/me", handlers.Me.Hello)
		auth.Put("/me", middlewares.NoBots, requests.Use[requests.UpdateMe](), middlewares.NewSimpleLimiter(3, 10*time.Minute, "Too many modifications, try again later"), handlers.Me.Update)
		auth.Get("/me/ip-check", handlers.Me.SwitchIPCheck)
		auth.Post("/me/profile-image", middlewares.NoBots, handlers.Me.UploadProfileImage)
		auth.Delete("/me", middlewares.NoBots, handlers.Me.Delete)
	}

	mobile := r.Group("/mobile", middlewares.NoBots)
	{
		mobile.Post("/register-push-notifications", requests.Use[requests.PushToken](), handlers.Mobile.RegisterPushNotification)
	}

	security := auth.Group("/security", middlewares.NoBots)
	{
		security.Get("/sessions", handlers.Security.Sessions)
		security.Delete("/sessions", handlers.Security.DestroySessions)
		security.Delete("/sessions/:id<int>", handlers.Security.DestroySession)
		security.Get("/ips", handlers.Security.IPs)
		security.Delete("/ips/:id<int>", handlers.Security.DestroyIP)
		security.Use(HandleNotFoundError)
	}

	users := auth.Group("/users", middlewares.NoBots)
	{
		users.Post("/findByEmail", requests.Use[requests.Email](), handlers.Users.FindByEmail)
		users.Use(HandleNotFoundError)
	}

	auth.Get("/chat/contacts", handlers.Chat.Contacts)
	auth.Get("/chat/user/:id<int>", middlewares.GroupAccess("id"), handlers.Group.GetInfo)
	auth.Get("/chat/private/:id<int>", handlers.Chat.PrivateMessage)
	auth.Post("/chat/search", requests.Use[requests.SearchContacts](), handlers.Chat.Search)

	auth.Post("/chat/group", requests.Use[requests.NewGroup](), handlers.Group.New)
	auth.Get("/chat/roles", handlers.Group.GetAllRoles)

	chat := auth.Group("/chat/:id<int>", middlewares.GroupAccess("id"))
	{
		chat.Get("/", handlers.Chat.Find)
		chat.Put("/", middlewares.NewRoleMiddleware(role.EditGroupName), requests.Use[requests.EditGroupName](), handlers.Group.EditName)
		chat.Get("/paginate", handlers.Chat.Messages)
		chat.Put("/invite", middlewares.NewRoleMiddleware(role.Creator), requests.Use[requests.CustomInvite](), middlewares.NewSimpleLimiter(3, 15*time.Minute, "too many requests"), handlers.Group.CustomInvite)
		chat.Put("/emoji", middlewares.NewRoleMiddleware(role.Admin), requests.Use[requests.Emoji](), handlers.Group.Emoji)
		chat.Get("/leave", handlers.Group.Leave)
		chat.Delete("/", middlewares.NewRoleMiddleware(role.Creator), handlers.Group.Delete)
		chat.Put("/roles", middlewares.NewRoleMiddleware(role.Admin), requests.Use[requests.ModifyRoles](), handlers.Group.ModifyRoles)
		chat.Get("/message/:messageID", handlers.Chat.FindMessage)
		chat.Post("/new-invite", middlewares.NewRoleMiddleware(role.InviteUser), requests.Use[requests.NewInvite](), middlewares.NewSimpleLimiter(3, 10*time.Minute, "Try again later before creating another"), handlers.Invite.Create)

		chat.Post("/file", middlewares.NewRoleMiddleware(role.AttachFile), handlers.Chat.UploadFile)
		chat.Post("/group-image", middlewares.NewRoleMiddleware(role.EditGroupImage), handlers.Group.UploadGroupImage)

		chat.Get("/users", handlers.Group.Users)
		chat.Get("/user_count", handlers.Group.UserCount)
		chat.Use(HandleNotFoundError)
	}

	auth.Get("/invite/:code", handlers.Invite.Use)

	// bot := r.Group("/bots", middlewares.Auth)
	{
		// bot.Get("/")
		// bot.Post("/")
		// bot.Post("/image")
		// bot.Put("/")
		// bot.Delete("/")
		// bot.Use(HandleNotFoundError)
	}

	auth.Use(HandleNotFoundError)
	r.Use(HandleNotFoundError)
}
