package redis

import (
	"BookShop/user_service/internal/config"
	"github.com/redis/go-redis/v9"
)

type RedisDb struct {
	client *redis.Client
}

func NewClient(cfg config.RedisConfig) (*RedisDb, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: "",
		DB:       0,
	})
	return &RedisDb{client: client}, nil
}
