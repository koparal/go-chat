package cache

import (
	"github.com/go-redis/redis/v8"
	"net"
)

type Config struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

func New(redisConf Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(redisConf.Host, redisConf.Port),
		Password: "",
		DB:       0,
	})
}
