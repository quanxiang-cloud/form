package models

import "gorm.io/gorm"

// PerGrant PerGrant
type PerGrant struct {
	PerGroupID string
	Owner      string
	OwnerName  string
	Types      int
}

type PerGrantQuery struct {
	PerGroupID string
	Permits    []string
	Owners     []string
}

type PerGrantRepo interface {
	BatchCreate(db *gorm.DB, perGrant ...*PerGrant) error
	Get(db *gorm.DB, permit, owner string) (*PerGrant, error)
	Find(db *gorm.DB, query *PerGrantQuery) ([]*PerGrant, error)
	Delete(db *gorm.DB, query *PerGrantQuery) error
}
