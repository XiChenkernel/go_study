package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"os"
	"strconv"
)

var RedisClient *redis.Client

func Redis() {
	db, _ := strconv.ParseUint(os.Getenv("REDIS_DB"), 10, 64)
	client := redis.NewClient(&redis.Options{
		Addr:       "127.0.0.1:6379",
		Password:   "",
		DB:         int(db),
		MaxRetries: 1,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic("can`t connect redis")
	}
	RedisClient = client
}
