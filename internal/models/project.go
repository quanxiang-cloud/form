package models

import "gorm.io/gorm"

type Project struct {
	ID          string
	Name        string
	Description string
	CreatedAt   int64
	CreatorID   string
	CreatorName string
}

type ProjectQuery struct {
	ID string
}

// ProjectRepo projectRepo.
type ProjectRepo interface {
	BatchCreate(db *gorm.DB, p ...*Project) error

	Get(db *gorm.DB, id string) (*Project, error)

	Delete(db *gorm.DB, query *ProjectQuery) error

	Update(db *gorm.DB, query *ProjectQuery, p *Project) error

	List(db *gorm.DB, query *ProjectQuery, page, size int) ([]*Project, int64, error)
}
