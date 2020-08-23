package service

import (
	"context"

	"github.com/go-redis/redis/v8"

	rc "go-kunpeng/config/redis"
)

func CreateRedisClient() (*redis.Client, error) {
	c := redis.NewClient(&redis.Options{
		Addr:     rc.Addr,
		Password: rc.Password,
		DB:       rc.DB,
	})

	var ctx = context.Background()
	_, err := c.Ping(ctx).Result()

	return c, err
}
