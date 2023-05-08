package dependencies

import (
	"context"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"ralts/internal/config"
	"time"
)

type CoreCacheInterface interface {
	Set(key string, value interface{}, expr time.Duration) error
	Get(key string) (string, error)
	Incr(key string) error
	Decr(key string) error
}

type Cache struct {
	Client *redis.Client
}

var ctx = context.Background()

func NewCache(cfg *config.Config) *Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisConn,
		DB:   0, // use default DB
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("unable to start Redis instance: %e", err)
	}

	return &Cache{
		Client: rdb,
	}
}

func (c *Cache) Set(key string, value interface{}, expr time.Duration) error {
	return c.Client.Set(ctx, key, value, expr).Err()
}

func (c *Cache) Get(key string) (string, error) {
	r, err := c.Client.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}

	return r, nil
}

func (c *Cache) Incr(key string) error {
	return c.Client.Incr(ctx, key).Err()
}

func (c *Cache) Decr(key string) error {
	return c.Client.Decr(ctx, key).Err()
}
