package models

import (
	"gorm.io/gorm"
)

const (
	// InitType 默认初始化的权限
	InitType PerType = 1
	// CreateType CreateType
	CreateType PerType = 2
)

// PerType PerType
type PerType int64

// PerGroup PerGroup
type PerGroup struct {
	ID          string
	AppID       string
	Name        string
	Description string
	CreatedAt   int64
	CreatorID   string
	CreatorName string
	Types       PerType
}

type PerGroupQuery struct {
	ID          string
	AppID       string
	DepID       string
	UserID      string
	Name        string
	PerGroupIDs []string
	Types       PerType
}

// PerGroupRepo PerGroupRepo
type PerGroupRepo interface {
	BatchCreate(db *gorm.DB, permissionGroup ...*PerGroup) error
	Get(db *gorm.DB, id string) (*PerGroup, error)
	Find(db *gorm.DB, query *PerGroupQuery) ([]*PerGroup, error)
	Update(db *gorm.DB, id string, perGroup *PerGroup) error
	Delete(db *gorm.DB, group *PerGroupQuery) error
}
