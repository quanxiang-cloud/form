package redis

import (
	"context"
	"encoding/json"
	"github.com/quanxiang-cloud/form/internal/models"
	"time"

	"github.com/go-redis/redis/v8"
)

type limitRepo struct {
	c *redis.ClusterClient
}

func (p *limitRepo) CreatePermit(ctx context.Context, roleID string, limits []*models.Limits) error {
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

func (p *limitRepo) GetPermit(ctx context.Context, roleID, path string) (*models.Limits, error) {
	key := p.PerKey() + roleID + ":" + path
	entityByte, err := p.c.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	entity := new(models.Limits)
	err = json.Unmarshal(entityByte, entity)
	return entity, err
}

func (p *limitRepo) DeletePermit(ctx context.Context, roleID, path string) error {

	return p.c.Del(ctx, p.PerKey()+roleID+":"+path).Err()
}

func (p *limitRepo) CreatePerMatch(ctx context.Context, match *models.PermitMatch) error {
	return p.c.HSet(ctx, p.PerMatchKey()+
		match.AppID, match.UserID, match.PermitID).Err()
}

func (p *limitRepo) GetPerMatch(ctx context.Context, appID, userID string) (*models.PermitMatch, error) {
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
	resp.PermitID = result.Val()
	return resp, nil
}

func (p *limitRepo) DeletePerMatch(ctx context.Context, appID string) error {
	return p.c.Del(ctx, p.PerMatchKey()+appID).Err()
}

func (p *limitRepo) Lock(ctx context.Context, key string, val interface{}, ttl time.Duration) (bool, error) {
	return p.c.SetNX(ctx, p.LockKey()+key, val, ttl).Result()
}

func (p *limitRepo) UnLock(ctx context.Context, key string) error {
	return p.c.Del(ctx, p.LockKey()+key).Err()
}

func (p *limitRepo) PerMatchExpire(ctx context.Context, key string, ttl time.Duration) error {
	return p.c.Expire(ctx, p.PerMatchKey()+key, ttl).Err()
}

func (p *limitRepo) PermitExpire(ctx context.Context, key string, ttl time.Duration) error {
	return p.c.Expire(ctx, p.PerKey()+key, ttl).Err()
}

//NewLimitRepo NewLimitRepo
func NewLimitRepo(c *redis.ClusterClient) models.LimitsRepo {
	return &limitRepo{
		c: c,
	}
}
func (p *limitRepo) PerKey() string {
	return redisKey + ":perInfo:"
}
func (p *limitRepo) PerMatchKey() string {
	return redisKey + ":perMatch:"
}
func (p *limitRepo) LockKey() string {
	return redisKey + ":lock:"
}
