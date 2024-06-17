package routes

import (
	"github.com/wizzldev/chat/app/handlers"
	"github.com/wizzldev/chat/app/services"
)

func MustInitApplication() {
	store, err := services.NewStorage()
	if err != nil {
		panic(err)
	}

	handlers.Chat.Init(store, services.NewWSCache())
	handlers.Files.Init(store)
	handlers.Group.Init(store)
	handlers.Me.Init(store)
}
