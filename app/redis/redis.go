package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/largezhou/lz_tools_backend/app/config"
)

var cfg = config.Config.Redis
var Client *redis.Client

func init() {
	Client = redis.NewClient(&redis.Options{
		Addr:     cfg.Host,
		Password: cfg.Password,
		DB:       cfg.Db,
	})
}
