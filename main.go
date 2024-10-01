package main

import (
	"flag"
	"log"

	"github.com/wizzldev/chat/database"
	"github.com/wizzldev/chat/pkg/configs"
	"github.com/wizzldev/chat/routes"
)

func main() {
	envFile := flag.String("env", ".env", "dotenv file to load")
	flag.Parse()

	err := configs.LoadEnv(*envFile)
	if err != nil {
		log.Fatal(err)
	}

	database.MustConnect()

	app := routes.NewApp()

	log.Fatal(app.Listen(configs.Env.ServerPort))
}
