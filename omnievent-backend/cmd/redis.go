package cmd

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitRedis() error {
	cfg := GetConfig()
	redisClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		DB:   cfg.Redis.DB,
	})

	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	return err
}

func GetRedis() *redis.Client {
	return redisClient
}

func CloseRedis() error {
	if redisClient != nil {
		return redisClient.Close()
	}
	return nil
}
