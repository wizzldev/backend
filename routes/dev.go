package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizzldev/chat/app/handlers"
	"github.com/wizzldev/chat/pkg/middlewares"
)

func RegisterDev(r fiber.Router) {
	auth := r.Group("/", middlewares.Auth)
	auth.Get("/applications", handlers.Developers.GetApplications)
}
