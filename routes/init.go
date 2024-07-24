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
	wsCache := services.NewWSCache()

	handlers.Chat.Init(store, wsCache)
	handlers.Files.Init(store)
	handlers.Group.Init(store, services.NewWSCache())
	handlers.Me.Init(store)
	handlers.Invite.Init(wsCache)
}
