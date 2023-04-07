package mysql

import (
	"github.com/quanxiang-cloud/form/internal/models"
	"gorm.io/gorm"
)

type projectUserRepo struct {
}

func (projectUser *projectUserRepo) BatchCreate(db *gorm.DB, p ...*models.ProjectUser) error {
	return db.Table(projectUser.TableName()).CreateInBatches(p, len(p)).Error
}

func (projectUser *projectUserRepo) Delete(db *gorm.DB, query *models.ProjectUserQuery) error {
	resp := make([]models.Project, 0)
	ql := db.Table(projectUser.TableName())
	if query.ProjectID != "" {
		ql = ql.Where("project_id = ? ", query.ProjectID)
	}
	if len(query.UserIDs) != 0 {
		ql = ql.Where("user_id in ?", query.UserIDs)
	}
	return ql.Delete(resp).Error
}

func (projectUser *projectUserRepo) List(db *gorm.DB, query *models.ProjectUserQuery, page, size int) ([]*models.ProjectUser, int64, error) {
	db = db.Table(projectUser.TableName())
	if query.ProjectID != "" {
		db = db.Where("project_id = ?", query.ProjectID)
	}
	if query.UserID != "" {
		db = db.Where("user_id = ?", query.UserID)
	}
	var (
		count    int64
		projects []*models.ProjectUser
	)
	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Offset((page - 1) * size).Limit(size).Find(&projects).Error
	if err != nil {
		return nil, 0, err
	}

	return projects, count, nil
}

func (projectUser *projectUserRepo) TableName() string {
	return "project_user"
}

func NewProjectUserRepo() models.ProjectUserRepo {
	return &projectUserRepo{}
}
