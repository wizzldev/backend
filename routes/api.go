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
		auth.Post("/me/profile-image", middlewares.NoBots, handlers.Me.UploadProfileImage)
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

	chat := auth.Group("/chat")
	{
		chat.Get("/contacts", handlers.Chat.Contacts)
		chat.Get("/user/:id<int>", middlewares.GroupAccess("id"), handlers.Group.GetInfo)
		chat.Get("/private/:id<int>", handlers.Chat.PrivateMessage)
		chat.Post("/search", requests.Use[requests.SearchContacts](), handlers.Chat.Search)
	}

	msg := chat.Group("/:id<int>", middlewares.GroupAccess("id"))
	{
		msg.Get("/", handlers.Chat.Find)
		msg.Put("/", middlewares.NewRoleMiddleware(role.EditGroupName), requests.Use[requests.EditGroupName](), handlers.Group.EditName)
		msg.Put("/roles", middlewares.NewRoleMiddleware(role.Admin), requests.Use[requests.ModifyRoles](), handlers.Group.ModifyRoles)
		msg.Post("/group-image", middlewares.NewRoleMiddleware(role.EditGroupImage), handlers.Group.UploadGroupImage)
		msg.Get("/paginate", handlers.Chat.Messages)
		msg.Post("/file", middlewares.NewRoleMiddleware(role.AttachFile), handlers.Chat.UploadFile)
		msg.Get("/message/:messageID", handlers.Chat.FindMessage)
		msg.Use(HandleNotFoundError)
	}

	// bot := r.Group("/bots", middlewares.Auth)
	{
		// bot.Get("/")
		// bot.Post("/")
		// bot.Post("/image")
		// bot.Put("/")
		// bot.Delete("/")
		// bot.Use(HandleNotFoundError)
	}

	chat.Post("/group", requests.Use[requests.NewGroup](), handlers.Group.New)
	chat.Get("/roles", handlers.Group.GetAllRoles)
	_ = chat.Group("/group/:groupId")

	auth.Use(HandleNotFoundError)
	chat.Use(HandleNotFoundError)
	r.Use(HandleNotFoundError)
}
