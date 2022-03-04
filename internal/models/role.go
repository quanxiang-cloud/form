package models

import "gorm.io/gorm"

const (
	// InitType 默认初始化的权限
	InitType RoleType = 1
	// CreateType CreateType
	CreateType RoleType = 2
)

// RoleType RoleType
type RoleType int64

type Role struct {
	ID          string
	AppID       string
	Name        string
	Description string
	CreatedAt   int64
	CreatorID   string
	CreatorName string
	Types       RoleType
}

type RoleQuery struct {
	ID      string
	AppID   string
	Name    string
	RoleIDS []string
	Types   RoleType
}

// RoleRepo RoleRepo
type RoleRepo interface {
	BatchCreate(db *gorm.DB, role ...*Role) error
	Get(db *gorm.DB, id string) (*Role, error)
	Find(db *gorm.DB, query *RoleQuery) ([]*Role, error)
	Update(db *gorm.DB, id string, role *Role) error
	Delete(db *gorm.DB, query *RoleQuery) error
}
