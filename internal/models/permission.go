package models

import (
	"context"
	"time"
)

// Permission Permission
type Permission struct {
	PerGroupID string
	FormID     string
	AppID      string
	Name       string
	Authority  int64
	Conditions map[string]Query
	Filter     map[string]interface{}
	Type       PerType
}

// PermissionMatch PermissionMatch
type PermissionMatch struct {
	PerGroupID string
	UserID     string
	AppID      string
}

// PermissionRepo 操作redis 的 接口
type PermissionRepo interface {
	Get(ctx context.Context, perGroupID, formID string) (*Permission, error)

	Create(ctx context.Context, permission *Permission, ttl time.Duration) error
	// Delete 删除
	Delete(ctx context.Context, permission *Permission) error
	// CreatePerMatch 根据 ID 和 FormID 得到 权限组id
	CreatePerMatch(ctx context.Context, match *PermissionMatch) error

	GetPerMatch(ctx context.Context, userID, appID string) (*PermissionMatch, error)
	// DeletePerMatch 删除perMatch
	DeletePerMatch(ctx context.Context, appID string) error
	// Lock 设置分布式锁
	Lock(ctx context.Context, key string, val interface{}, ttl time.Duration) (bool, error)
	// UnLock 解除分布式锁
	UnLock(ctx context.Context, key string) error
	// PerMatchExpire 给某个键设置过期时间
	PerMatchExpire(ctx context.Context, key string, ttl time.Duration) error
}
