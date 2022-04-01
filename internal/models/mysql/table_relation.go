package mysql

import (
	"gorm.io/gorm"

	"github.com/quanxiang-cloud/form/internal/models"
)

type tableRelationRepo struct{}

func (t *tableRelationRepo) Get(db *gorm.DB, tableID, fieldName string) (*models.TableRelation, error) {
	tableRelation := new(models.TableRelation)
	err := db.Table(t.TableName()).Where("table_id = ? and  field_name = ? ", tableID, fieldName).Find(tableRelation).Error
	if err != nil {
		return nil, err
	}
	return tableRelation, nil
}

func (t *tableRelationRepo) BatchCreate(db *gorm.DB, table ...*models.TableRelation) error {
	return db.Table(t.TableName()).CreateInBatches(table, len(table)).Error
}

func (t *tableRelationRepo) Find(db *gorm.DB, query *models.TableRelationQuery) ([]*models.TableRelation, error) {
	return nil, nil
}

func (t *tableRelationRepo) Update(db *gorm.DB, tableID, fieldName string, table *models.TableRelation) error {
	setMap := make(map[string]interface{})
	if table.Filter != nil {
		setMap["filter"] = table.Filter
	}
	if table.SubTableID != "" {
		setMap["sub_table_id"] = table.SubTableID
	}
	if table.SubTableType != "" {
		setMap["sub_table_type"] = table.SubTableType
	}
	return db.Table(t.TableName()).Where("table_id = ? and field_name = ? ", tableID, fieldName).Updates(
		setMap).Error
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
