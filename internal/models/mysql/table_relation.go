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
	resp := make([]models.TableRelation, 0)
	ql := db.Table(t.TableName())
	if query.TableID != "" {
		ql = ql.Where("table_id = ? ", query.TableID)
	}
	if query.AppID != "" {
		ql = ql.Where("app_id = ?", query.AppID)
	}
	return ql.Delete(resp).Error
}
func (t *tableRelationRepo) List(db *gorm.DB, query *models.TableRelationQuery, page, size int) ([]*models.TableRelation, int64, error) {
	db = db.Table(t.TableName())
	if query.AppID != "" {
		db = db.Where("app_id = ?", query.AppID)
	}
	if query.SubTableID != "" {
		db = db.Where("sub_table_id = ?", query.SubTableID)
	}
	if query.SubTableType != "" {
		db = db.Where("sub_table_type = ?", query.SubTableType)
	}
	if query.FieldName != "" {
		db = db.Where("field_name = ?", query.FieldName)
	}
	if query.TableID != "" {
		db = db.Where("table_id = ?", query.TableID)
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

func NewTableRelationRepo() models.TableRelationRepo {
	return &tableRelationRepo{}
}
