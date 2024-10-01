package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/app/handlers"
	"github.com/wizzldev/chat/app/requests"
	"github.com/wizzldev/chat/pkg/middlewares"
)

func RegisterDev(r fiber.Router) {
	auth := r.Group("/", middlewares.Auth)

	auth.Get("/applications", handlers.Developers.GetApplications).Name("apps")
	auth.Post("/applications", requests.Use[requests.NewBot](), handlers.Developers.CreateApplication).Name("apps")
	auth.Patch("/applications/:id<int>", handlers.Developers.RegenerateApplicationToken).Name("apps")

	auth.Post("/invite", requests.Use[requests.ApplicationInvite](), handlers.Group.InviteApplication).Name("apps.invite")
}
