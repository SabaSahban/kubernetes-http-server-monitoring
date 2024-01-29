package storage

import (
	"cloud-final/config"
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type Storage struct {
	timeout     int
	redisClient *redis.Client
}

func New(cfg *config.Config) *Storage {
	rc := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host,
		Password: cfg.Redis.Password,
		DB:       0,
	})

	logrus.Info("redis config host", cfg.Redis.Host)

	return &Storage{
		timeout:     cfg.Redis.Timeout,
		redisClient: rc,
	}
}

// Read key from the database.
func (s *Storage) Read(key string) (string, error) {
	ctx := context.Background()

	return s.redisClient.Get(ctx, key).Result()
}

// Write a set into the database.
func (s *Storage) Write(key string, value interface{}) error {
	ctx := context.Background()

	// Serialize the JSON data
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// Store the JSON data in Redis
	return s.redisClient.Set(ctx, key, jsonValue, time.Duration(s.timeout)*time.Minute).Err()
}

func (s *Storage) Ping(ctx context.Context) error {
	return s.redisClient.Ping(ctx).Err()
}
