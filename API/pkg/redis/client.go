package redis

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/AmadoMuerte/BirthdayWish/API/pkg/config"
	"github.com/redis/go-redis/v9"
)

var (
	instance *RDB
	once     sync.Once
)

type RDB struct {
	Client *redis.Client
}

func GetInstance(cfg *config.Config) (*RDB, error) {
	var initErr error
	once.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Path, cfg.Redis.Port),
			Username: cfg.Redis.User,
			Password: cfg.Redis.Pass,
			DB:       0,
		})

		if _, err := client.Ping(context.Background()).Result(); err != nil {
			initErr = fmt.Errorf("failed to connect to Redis: %w", err)
			return
		}

		instance = &RDB{Client: client}
		slog.Info("Redis connection established")
	})

	return instance, initErr
}
