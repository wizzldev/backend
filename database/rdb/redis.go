package rdb

import (
	"fmt"
	"github.com/gofiber/storage/redis/v3"
	client "github.com/redis/go-redis/v9"
	"github.com/wizzldev/chat/pkg/configs"
)

var Redis *redis.Storage

var RedisClient client.UniversalClient

func MustConnect() {
	fmt.Printf("Connecting to redis via %s:%d", configs.Env.Redis.Host, configs.Env.Redis.Port)

	Redis = redis.New(redis.Config{
		Host:     configs.Env.Redis.Host,
		Port:     configs.Env.Redis.Port,
		Username: configs.Env.Redis.User,
		Password: configs.Env.Redis.Password,
		Database: configs.Env.Redis.DB,
	})

	RedisClient = Redis.Conn()
}
