package storage

import (
	"cloud-final/config"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type Storage struct {
	timeout     int
	redisClient *redis.Client
}

type ServerData struct {
	ID           string    `json:"id"`
	Address      string    `json:"address"`
	SuccessCount int       `json:"success"`
	FailureCount int       `json:"failure"`
	LastFailure  time.Time `json:"last_failure"`
	CreatedAt    time.Time `json:"created_at"`
}

func New(cfg *config.Config) *Storage {
	rc := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host,
		Password: cfg.Redis.Password,
		DB:       0,
	})

	logrus.Info("Redis config host", cfg.Redis.Host)

	return &Storage{
		timeout:     cfg.Redis.Timeout,
		redisClient: rc,
	}
}

// AddServer adds a new server to the storage.
func (s *Storage) AddServer(data ServerData) error {
	ctx := context.Background()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("server:%s", data.ID)

	return s.redisClient.Set(ctx, key, jsonData, 0).Err()
}

func (s *Storage) GetServer(id string) (*ServerData, error) {
	ctx := context.Background()

	result, err := s.redisClient.Get(ctx, fmt.Sprintf("server:%s", id)).Result()
	if err != nil {
		return nil, err
	}

	var data ServerData
	err = json.Unmarshal([]byte(result), &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (s *Storage) GetAllServers() ([]ServerData, error) {
	ctx := context.Background()

	keys, err := s.redisClient.Keys(ctx, "server:*").Result()
	if err != nil {
		return nil, err
	}

	var servers []ServerData

	for _, key := range keys {
		result, err := s.redisClient.Get(ctx, key).Result()
		if err != nil {
			return nil, err
		}

		var server ServerData
		err = json.Unmarshal([]byte(result), &server)
		if err != nil {
			return nil, err
		}

		servers = append(servers, server)
	}

	return servers, nil
}

func (s *Storage) UpdateServer(data ServerData) error {
	ctx := context.Background()
	key := fmt.Sprintf("server:%s", data.ID)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return s.redisClient.Set(ctx, key, jsonData, 0).Err()
}

func (s *Storage) Ping(ctx context.Context) error {
	return s.redisClient.Ping(ctx).Err()
}
