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
		ErrorHandler: routes.ErrorHandler,
	})

	app.Use(recover.New())
	app.Use(middlewares.CORS())

	if configs.Env.Debug {
		app.Use(logger.New())
	}

	routes.RegisterStorage(app.Group("/storage"))
	routes.RegisterAPI(app.Group("/api"))
	routes.WS(app)

	app.Use(routes.HandleNotFoundError)

	log.Fatal(app.Listen(configs.Env.ServerPort))
}
