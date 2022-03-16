package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/quanxiang-cloud/form/internal/models"
)

type LimitRepo struct {
	c *redis.ClusterClient
}

// NewLimitRepo NewLimitRepo
func NewLimitRepo(c *redis.ClusterClient) *LimitRepo {
	return &LimitRepo{
		c: c,
	}
}

func (p *LimitRepo) ExistsKey(ctx context.Context, key string) bool {
	exists := p.c.Exists(ctx, key)
	return exists.Val() > 0
}

func (p *LimitRepo) CreatePermit(ctx context.Context, roleID string, limits ...*models.Limits) error {
	key := p.PerKey() + roleID
	for _, value := range limits {
		entityJSON, err := json.Marshal(value)
		if err != nil {
			return nil
		}
		p.c.HSet(ctx, key+value.Path, entityJSON)
	}
	return nil
}

func (p *LimitRepo) GetPermit(ctx context.Context, roleID, path string) (*models.Limits, error) {
	result := p.c.HGet(ctx, p.PerKey()+roleID, path)
	if result.Err() == redis.Nil {
		return nil, nil
	}
	if result.Err() != nil {
		return nil, result.Err()
	}
	bytes, err := result.Bytes()
	if err != nil {
		return nil, err
	}
	limits := new(models.Limits)
	err = json.Unmarshal(bytes, limits)
	if err != nil {
		return nil, err
	}
	return limits, nil
}

func (p *LimitRepo) DeletePermit(ctx context.Context, roleID string) error {
	return p.c.Del(ctx, p.PerKey()+roleID).Err()
}

func (p *LimitRepo) CreatePerMatch(ctx context.Context, match *models.PermitMatch) error {
	return p.c.HSet(ctx, p.PerMatchKey()+
		match.AppID, match.UserID, match.RoleID).Err()
}

func (p *LimitRepo) GetPerMatch(ctx context.Context, appID, userID string) (*models.PermitMatch, error) {
	result := p.c.HGet(ctx, p.PerMatchKey()+appID, userID)
	if result.Err() == redis.Nil {
		return nil, nil
	}
	if result.Err() != nil {
		return nil, result.Err()
	}
	resp := &models.PermitMatch{
		UserID: userID,
		AppID:  appID,
	}
	resp.RoleID = result.Val()
	return resp, nil
}

func (p *LimitRepo) DeletePerMatch(ctx context.Context, appID string) error {
	return p.c.Del(ctx, p.PerMatchKey()+appID).Err()
}

func (p *LimitRepo) Lock(ctx context.Context, key string, val interface{}, ttl time.Duration) (bool, error) {
	return p.c.SetNX(ctx, p.LockKey()+key, val, ttl).Result()
}

func (p *LimitRepo) UnLock(ctx context.Context, key string) error {
	return p.c.Del(ctx, p.LockKey()+key).Err()
}

func (p *LimitRepo) PerMatchExpire(ctx context.Context, key string, ttl time.Duration) error {
	return p.c.Expire(ctx, p.PerMatchKey()+key, ttl).Err()
}

func (p *LimitRepo) PermitExpire(ctx context.Context, key string, ttl time.Duration) error {
	return p.c.Expire(ctx, p.PerKey()+key, ttl).Err()
}

func (p *LimitRepo) PerKey() string {
	return redisKey + ":perInfo:"
}

func (p *LimitRepo) PerMatchKey() string {
	return redisKey + ":perMatch:"
}

func (p *LimitRepo) LockKey() string {
	return redisKey + ":lock:"
}
