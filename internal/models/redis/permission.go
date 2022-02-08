package redis

import (
	"context"
	"encoding/json"
	"github.com/quanxiang-cloud/form/internal/models"
	"time"

	"github.com/go-redis/redis/v8"
)

type permissionRepo struct {
	c *redis.ClusterClient
}

//NewPermissionRepo NewPermissionRepo
func NewPermissionRepo(c *redis.ClusterClient) models.PermissionRepo {
	return &permissionRepo{
		c: c,
	}
}
func (p *permissionRepo) PerKey() string {
	return redisKey + ":perInfo:"
}
func (p *permissionRepo) PerMatchKey() string {
	return redisKey + ":perMatch:"
}
func (p *permissionRepo) LockKey() string {
	return redisKey + ":lock:"
}

func (p *permissionRepo) Get(ctx context.Context, perGroupID, formID string) (*models.Permission, error) {
	key := p.PerKey() + perGroupID + ":" + formID
	entityByte, err := p.c.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	entity := new(models.Permission)
	err = json.Unmarshal(entityByte, entity)
	return entity, err

}

func (p *permissionRepo) Create(ctx context.Context, permission *models.Permission, ttl time.Duration) error {
	entityJSON, err := json.Marshal(permission)
	if err != nil {
		return err
	}
	key := p.PerKey() + permission.PerGroupID + ":" + permission.FormID
	return p.c.Set(ctx, key, entityJSON, ttl).Err()
}

func (p *permissionRepo) Delete(ctx context.Context, permission *models.Permission) error {
	key := p.PerKey() + permission.PerGroupID + ":" + permission.FormID
	return p.c.Del(ctx, key).Err()
}

func (p *permissionRepo) CreatePerMatch(ctx context.Context, match *models.PermissionMatch) error {
	return p.c.HSet(ctx, p.PerMatchKey()+
		match.AppID, match.UserID, match.PerGroupID).Err()

}

func (p *permissionRepo) GetPerMatch(ctx context.Context, userID, appID string) (*models.PermissionMatch, error) {

	result := p.c.HGet(ctx, p.PerMatchKey()+appID, userID)
	if result.Err() == redis.Nil {
		return nil, nil
	}
	if result.Err() != nil {
		return nil, result.Err()
	}
	resp := &models.PermissionMatch{
		UserID: userID,
		AppID:  appID,
	}
	resp.PerGroupID = result.Val()
	return resp, nil
}
func (p *permissionRepo) DeletePerMatch(ctx context.Context, appID string) error {
	return p.c.Del(ctx, p.PerMatchKey()+appID).Err()
}

func (p *permissionRepo) Lock(ctx context.Context, key string, val interface{}, ttl time.Duration) (bool, error) {
	return p.c.SetNX(ctx, p.LockKey()+key, val, ttl).Result()
}

func (p *permissionRepo) UnLock(ctx context.Context, key string) error {
	return p.c.Del(ctx, p.LockKey()+key).Err()
}

func (p *permissionRepo) PerMatchExpire(ctx context.Context, key string, ttl time.Duration) error {
	return p.c.Expire(ctx, p.PerMatchKey()+key, ttl).Err()
}
