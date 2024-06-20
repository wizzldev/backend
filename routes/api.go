package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/app/handlers"
	"github.com/wizzldev/chat/app/requests"
	"github.com/wizzldev/chat/pkg/middlewares"
)

func RegisterAPI(r fiber.Router) {
	{
		r.Post("/login", requests.Use[requests.Login](), middlewares.NewAuthLimiter(), handlers.Auth.Login)
		r.Post("/register", requests.Use[requests.Register](), handlers.Auth.Register)
		r.Get("/logout", handlers.Auth.Logout)
		r.Post("/request-new-password", requests.Use[requests.NewPassword](), handlers.Auth.RequestNewPassword)
		r.Get("/set-new-password/:token", handlers.Auth.IsPasswordResetExists)
		r.Post("/set-new-password/:token", requests.Use[requests.SetNewPassword](), handlers.Auth.SetNewPassword)
		r.Get("/verify-email/:token", handlers.Auth.VerifyEmail)
		r.Post("/request-new-email-verification", requests.Use[requests.Email](), handlers.Auth.RequestNewEmailVerification)
	}

	auth := r.Group("/", middlewares.Auth)
	{
		auth.Get("/me", handlers.Me.Hello)
		auth.Post("/me/profile-image", handlers.Me.UploadProfileImage)
	}

	users := auth.Group("/users")
	{
		users.Post("/findByEmail", requests.Use[requests.Email](), handlers.Users.FindByEmail)
	}

	chat := auth.Group("/chat")
	{
		chat.Get("/contacts", handlers.Chat.Contacts)
		chat.Get("/user/:id<int>", handlers.Group.GetInfo)
		chat.Get("/private/:id<int>", handlers.Chat.PrivateMessage)
		chat.Post("/search", requests.Use[requests.SearchContacts](), handlers.Chat.Search)
	}

	msg := chat.Group("/:id<int>", middlewares.GroupAccess("id"))
	{
		msg.Get("/", handlers.Chat.Find)
		msg.Get("/paginate", handlers.Chat.Messages)
		msg.Post("/file", handlers.Chat.UploadFile)
	}

	chat.Post("/group", requests.Use[requests.NewGroup](), handlers.Group.New)
	chat.Get("/roles", handlers.Group.GetAllRoles)
	_ = chat.Group("/group/:groupId")

	auth.Use(HandleNotFoundError)
	users.Use(HandleNotFoundError)
	chat.Use(HandleNotFoundError)
	msg.Use(HandleNotFoundError)
	r.Use(HandleNotFoundError)
}
