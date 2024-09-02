package rdb

import (
	"github.com/gofiber/storage/redis/v3"
	"os"
)

var Redis = redis.New(redis.Config{
	URL: getURL(),
})

func getURL() string {
	url := os.Getenv("REDIS_URL")
	if url == "" {
		return "redis://:@127.0.0.1:6379/0"
	}
	return url
}
