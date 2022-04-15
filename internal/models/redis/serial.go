package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/quanxiang-cloud/form/internal/models"
)

type serialRepo struct {
	c *redis.ClusterClient
}

// NewSerialRepo NewSerialRepo
func NewSerialRepo(c *redis.ClusterClient) models.SerialRepo {
	return &serialRepo{
		c: c,
	}
}

func (s *serialRepo) Key() string {
	return redisSerialKey
}

func (s *serialRepo) Create(ctx context.Context, appID, tableID, fieldID string, values map[string]interface{}) error {
	key := s.Key() + appID + ":" + tableID + ":" + fieldID
	return s.c.HSet(ctx, key, values).Err()
}

func (s *serialRepo) Get(ctx context.Context, appID, tableID, fieldID, field string) string {
	key := s.Key() + appID + ":" + tableID + ":" + fieldID
	return s.c.HGet(ctx, key, field).Val()
}

func (s *serialRepo) GetAll(ctx context.Context, appID, tableID, fieldID string) map[string]string {
	key := s.Key() + appID + ":" + tableID + ":" + fieldID
	return s.c.HGetAll(ctx, key).Val()
}
