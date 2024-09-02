package rdb

import (
	"github.com/gofiber/storage/redis/v3"
	client "github.com/redis/go-redis/v9"
	"github.com/wizzldev/chat/pkg/configs"
)

var Redis *redis.Storage

var RedisClient client.UniversalClient

func MustConnect() {
	Redis = redis.New(redis.Config{
		Host:     configs.Env.Redis.Host,
		Port:     configs.Env.Redis.Port,
		Username: configs.Env.Redis.User,
		Password: configs.Env.Redis.Password,
		Database: configs.Env.Redis.DB,
	})

	RedisClient = Redis.Conn()
}
