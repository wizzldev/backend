package routes

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/wizzldev/chat/pkg/configs"
	"github.com/wizzldev/chat/pkg/middlewares"
)

func NewApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler:       ErrorHandler,
		Prefork:            !configs.Env.Debug,
		ServerHeader:       "Wizzl",
		AppName:            "Wizzl v1.0.0",
		ProxyHeader:        fiber.HeaderXForwardedFor,
		EnableIPValidation: true,
	})

	if !configs.Env.Debug {
		app.Use(recover.New())
	} else {
		app.Use(logger.New())
	}

	app.Use(middlewares.CORS())

	MustInitApplication()
	RegisterRouteGetter(app)
	RegisterStorage(app.Group("/storage").Name("storage."))
	WS(app.Group("/ws"))
	RegisterBot(app.Group("/bot").Name("bot."))
	RegisterDev(app.Group("/developers").Name("devs."))
	RegisterAPI(app)
	app.Use(HandleNotFoundError)

	return app
}

func RegisterRouteGetter(r fiber.Router) {
	r.Get("/app/api-routes", func(c *fiber.Ctx) error {
		allRoute := c.App().GetRoutes()
		var routes []fiber.Route

		for _, route := range allRoute {
			if !strings.HasSuffix(route.Path, "/") && route.Name != "" {
				routes = append(routes, route)
			}
		}

		return c.JSON(routes)
	})
}
