package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/pkg/configs"
	"github.com/wizzldev/chat/pkg/middlewares"
	"github.com/wizzldev/chat/routes"
	"log"
)

func main() {
	err := configs.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}

	database.MustConnect()

	app := fiber.New(fiber.Config{
		ErrorHandler:       routes.ErrorHandler,
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

	routes.MustInitApplication()
	routes.RegisterStorage(app.Group("/storage"))
	routes.WS(app.Group("/ws"))
	routes.RegisterAPI(app)

	app.Use(routes.HandleNotFoundError)

	log.Fatal(app.Listen(configs.Env.ServerPort))
}
