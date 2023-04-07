package models

import "gorm.io/gorm"

type Project struct {
	ID           string
	Name         string
	Description  string
	CreatedAt    int64
	CreatorID    string
	CreatorName  string
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

type ProjectQuery struct {
	ID  string
	IDS []string
}

// ProjectRepo projectRepo.
type ProjectRepo interface {
	BatchCreate(db *gorm.DB, p ...*Project) error

	Get(db *gorm.DB, id string) (*Project, error)

	Delete(db *gorm.DB, query *ProjectQuery) error

	Update(db *gorm.DB, query *ProjectQuery, p *Project) error

	List(db *gorm.DB, query *ProjectQuery, page, size int) ([]*Project, int64, error)
}
