package models

import "gorm.io/gorm"

type UserRole struct {
	ID string
	// app id
	UserID string
	// table id
	RoleID string

	AppID string
}

type UserRoleQuery struct {
	UserID string
	AppID  string
}

// UserRoleRepo UserRoleRepo.
type UserRoleRepo interface {
	BatchCreate(db *gorm.DB, userRole ...*UserRole) error
	Get(db *gorm.DB, userID, appID string) (*UserRole, error)
	Delete(db *gorm.DB, query *UserRoleQuery) error
	List(db *gorm.DB, query *UserRoleQuery, page, size int) ([]*UserRole, int64, error)
}