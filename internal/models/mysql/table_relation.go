package mysql

import (
	"gorm.io/gorm"

	"github.com/quanxiang-cloud/form/internal/models"
)

type tableRelationRepo struct{}

func (t *tableRelationRepo) BatchCreate(db *gorm.DB, table ...*models.TableRelation) error {
	return nil
}

func (t *tableRelationRepo) Find(db *gorm.DB, query *models.TableRelationQuery) ([]*models.TableRelation, error) {
	return nil, nil
}

func (t *tableRelationRepo) Update(db *gorm.DB, table *models.TableRelation) error {
	return nil
}

func (t *tableRelationRepo) Delete(db *gorm.DB, query *models.TableRelationQuery) error {
	return nil
}

func (t *tableRelationRepo) List(db *gorm.DB, query *models.TableRelationQuery, page, size int) ([]*models.TableRelation, int64, error) {
	db = db.Table(t.TableName())
	if query.AppID != "" {
		db = db.Where("app_id = ?", query.AppID)
	}
	var (
		count          int64
		tableRelations []*models.TableRelation
	)

	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Order("created_at desc").Offset((page - 1) * size).Limit(size).Find(&tableRelations).Error
	if err != nil {
		return nil, 0, err
	}

	return tableRelations, count, nil
}

func (t *tableRelationRepo) TableName() string {
	return "table_relation"
}

func NewTableRelation() models.TableRelationRepo {
	return &tableRelationRepo{}
}
