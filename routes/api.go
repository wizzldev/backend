package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/app/handlers"
	"github.com/wizzldev/chat/app/requests"
	"github.com/wizzldev/chat/pkg/middlewares"
)

func RegisterAPI(r fiber.Router) {
	{
		r.Post("/login", requests.Use[requests.Login](), handlers.Auth.Login)
		r.Post("/register", requests.Use[requests.Register](), handlers.Auth.Register)
		r.Get("/logout", handlers.Auth.Logout)
	}

	auth := r.Group("/", middlewares.Auth)
	{
		auth.Get("/me", handlers.Me.Hello)
		auth.Post("/me/profile-image", handlers.Me.UploadProfileImage)
	}

	chat := auth.Group("/chat")
	{
		chat.Get("/contacts", handlers.Chat.Contacts)
		chat.Get("/private/:id<int>", handlers.Chat.PrivateMessage)
		chat.Post("/search", requests.Use[requests.SearchContacts](), handlers.Chat.Search)
	}

	msg := chat.Group("/:id<int>", middlewares.GroupAccess("id"))
	{
		msg.Get("/", handlers.Chat.Find)
	}

	chat.Post("/group", handlers.Group.New)
	_ = chat.Group("/group/:groupId")
}
