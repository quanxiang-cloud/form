package mysql

import (
	"github.com/quanxiang-cloud/form/internal/models"
	"gorm.io/gorm"
)

type projectRepo struct{}

func (project *projectRepo) BatchCreate(db *gorm.DB, p ...*models.Project) error {
	return db.Table(project.TableName()).CreateInBatches(p, len(p)).Error
}

func (project *projectRepo) Get(db *gorm.DB, id string) (*models.Project, error) {
	p := new(models.Project)
	err := db.Table(project.TableName()).Where("id = ? ", id).Find(p).Error
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (project *projectRepo) Delete(db *gorm.DB, query *models.ProjectQuery) error {
	resp := make([]models.Project, 0)
	ql := db.Table(project.TableName())
	if query.ID != "" {
		ql = ql.Where("id = ? ", query.ID)
	}
	return ql.Delete(resp).Error
}

func (project *projectRepo) Update(db *gorm.DB, query *models.ProjectQuery, p *models.Project) error {
	return nil
}

func (project *projectRepo) List(db *gorm.DB, query *models.ProjectQuery, page, size int) ([]*models.Project, int64, error) {
	db = db.Table(project.TableName())
	//if query.AppID != "" {
	//	db = db.Where("app_id = ?", query.AppID)
	//}
	var (
		count    int64
		projects []*models.Project
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

func (project *projectRepo) TableName() string {
	return "project"
}

func NewProjectRepo() models.ProjectRepo {
	return &projectRepo{}
}
