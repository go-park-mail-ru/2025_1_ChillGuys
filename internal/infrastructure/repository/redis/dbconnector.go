package redis

import (
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"github.com/redis/go-redis/v9"
)

// GetRedisOptions возвращает конфигурацию подключения к Redis
func GetRedisOptions(conf *config.RedisConfig) *redis.Options {
	return &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       conf.DB,
	}
}
