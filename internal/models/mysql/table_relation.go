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

func (t *tableRelationRepo) TableName() string {
	return "table_relation"
}

func NewTableRelation() models.TableRelationRepo {
	return &tableRelationRepo{}
}
