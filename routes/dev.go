package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/app/handlers"
	"github.com/wizzldev/chat/app/requests"
	"github.com/wizzldev/chat/pkg/middlewares"
)

func RegisterDev(r fiber.Router) {
	auth := r.Group("/", middlewares.Auth)

	auth.Get("/applications", handlers.Developers.GetApplications)
	auth.Post("/applications", requests.Use[requests.NewBot](), handlers.Developers.CreateApplication)
	auth.Patch("/applications/:id<int>", handlers.Developers.RegenerateApplicationToken)

	auth.Post("/invite", requests.Use[requests.ApplicationInvite](), handlers.Group.InviteApplication)
}
