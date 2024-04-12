package repository

import "github.com/redis/go-redis/v9"

var ClientRedis *redis.Client

const (
	ttlSeconds = 300
)

func InitRedisClient(addr, password string, db int) {
	ClientRedis = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}
