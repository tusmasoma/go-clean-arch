package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	"github.com/tusmasoma/go-clean-arch/config"
)

const (
	dbPrefix = "REDIS_"
)

func NewRedisClient(ctx context.Context) *redis.Client {
	conf, err := config.NewCacheConfig(ctx, dbPrefix)
	if err != nil || conf == nil {
		log.Error("Failed to load cache config: %s\n", log.Ferror(err))
		return nil
	}

	client := redis.NewClient(&redis.Options{Addr: conf.Addr, Password: conf.Password, DB: conf.DB})

	_, err = client.Ping(ctx).Result()
	if err != nil {
		log.Critical("Failed to connect to Redis", log.Ferror(err))
		return nil
	}

	log.Info("Successfully connected to Redis")
	return client
}
