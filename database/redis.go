package database

import (
	"github.com/gofiber/storage/redis/v3"
	"github.com/wizzldev/chat/pkg/configs"
)

var Redis = redis.New(redis.Config{
	Host:     configs.Env.Redis.Host,
	Port:     configs.Env.Redis.Port,
	Username: configs.Env.Redis.User,
	Password: configs.Env.Redis.Password,
	Database: configs.Env.Redis.DB,
})

var RedisClient = Redis.Conn()
