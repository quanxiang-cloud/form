package models

import "gorm.io/gorm"

type ProjectUser struct {
	ID          string
	ProjectID   string
	ProjectName string
	UserID      string
	UserName    string
}

type ProjectUserQuery struct {
	ProjectID    string
	UserIDs      []string
	UserID       string
	SerialNumber string
	// 开始时间
	StartAt int64
	// 结束时间
	EndAt int64
	// 状态
	Status string
	// 备注
	Remark string
}

// ProjectUserRepo projectUserRepo.
type ProjectUserRepo interface {
	BatchCreate(db *gorm.DB, p ...*ProjectUser) error

	Delete(db *gorm.DB, p *ProjectUserQuery) error

	List(db *gorm.DB, query *ProjectUserQuery, page, size int) ([]*ProjectUser, int64, error)
}
