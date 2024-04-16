package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/app/handlers"
	"github.com/wizzldev/chat/app/requests"
	"github.com/wizzldev/chat/pkg/middlewares"
)

func RegisterAPI(r fiber.Router) {
	{
		r.Post("/login", requests.Use(&requests.Login{}), handlers.Auth.Login)
		r.Post("/register", requests.Use(&requests.Register{}), handlers.Auth.Register)
	}

	auth := r.Group("/", middlewares.Auth)
	{
		auth.Get("/me", handlers.Me.Hello)
		auth.Post("/me/profile-image", handlers.Me.UploadProfileImage)
	}

	chat := auth.Group("/chat")
	{
		chat.Get("/private/:id", handlers.Chat.PrivateMessage)
		chat.Post("/search", requests.Use(&requests.SearchContacts{}), handlers.Chat.Search)
	}

	chat.Post("/group", handlers.Group.New)
	_ = chat.Group("/group/:groupId")
}
