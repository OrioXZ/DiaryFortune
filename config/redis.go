package config

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	Rdb *redis.Client
	Ctx = context.Background()
)

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"), // e.g., "redis:6379"
		Password: "",                      // no password set
		DB:       0,                       // use default DB
	})

	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}
}
