package redis_service

import (
	"app/config"
	"app/pkg/logs"
	"app/storage"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	redis *redis.Client
	log   logs.LoggerInterface
	cfg   *config.RedisConfig
	code  storage.CodeCacheInterface
}

var (
	ErrKeyNotFound = fmt.Errorf("key_not_found")
)

func NewRedisStore(cfg *config.RedisConfig, log logs.LoggerInterface) storage.CacheInterface {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: os.Getenv("REDIS_PWD"),
		DB:       0,
		// Username: os.Getenv("REDIS_USER"),
	})
	return &RedisStore{
		redis: client,
		code:  NewCodeCache(client, log),
		log:   log,
		cfg:   cfg,
	}
}

func (rs *RedisStore) Code() storage.CodeCacheInterface {
	if rs.code == nil {
		rs.code = NewCodeCache(rs.redis, rs.log)
	}
	return rs.code
}
