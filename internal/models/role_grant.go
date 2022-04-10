package models

import "gorm.io/gorm"

type RoleGrant struct {
	ID        string
	RoleID    string
	AppID     string
	Owner     string
	OwnerName string
	Types     int
	CreatedAt int64
}

type RoleGrantQuery struct {
	RoleID  string
	RoleIDs []string
	Owners  []string
	AppID   string
}

// RoleRantRepo RoleRantRepo
type RoleRantRepo interface {
	BatchCreate(db *gorm.DB, roleGrant ...*RoleGrant) error
	Get(db *gorm.DB, id string) (*RoleGrant, error)
	Update(db *gorm.DB, id string, roleGrant *RoleGrant) error
	Delete(db *gorm.DB, query *RoleGrantQuery) error
	List(db *gorm.DB, query *RoleGrantQuery, page, size int) ([]*RoleGrant, int64, error)
}
