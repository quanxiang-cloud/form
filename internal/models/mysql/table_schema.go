package mysql

import (
	"github.com/quanxiang-cloud/form/internal/models"
	"gorm.io/gorm"
)

type tableSchemaRepo struct{}

func NewTableSchema() models.TableSchemeRepo {
	return &tableSchemaRepo{}
}

func (t *tableSchemaRepo) TableName() string {
	return "table_schema"
}

func (t *tableSchemaRepo) BatchCreate(db *gorm.DB, schema ...*models.TableSchema) error {
	return db.Table(t.TableName()).CreateInBatches(schema, len(schema)).Error
}

func (t *tableSchemaRepo) Get(db *gorm.DB, appID, tableID string) (*models.TableSchema, error) {
	permitForm := new(models.TableSchema)
	err := db.Table(t.TableName()).Where("app_id = ? and  table_id = ? ", appID, tableID).Find(permitForm).Error
	if err != nil {
		return nil, err
	}
	return permitForm, nil
}

func (t *tableSchemaRepo) Find(db *gorm.DB, query *models.TableSchemaQuery, size int64, page int64) ([]*models.TableSchema, int64, error) {
	return nil, 0, nil
}

func (t *tableSchemaRepo) Delete(db *gorm.DB, query *models.TableSchemaQuery) error {
	return nil
}

func (t *tableSchemaRepo) Update(db *gorm.DB, appID, tableID string, tableSchema *models.TableSchema) error {
	setMap := make(map[string]interface{})
	if tableSchema.Schema != nil {
		setMap["schema"] = tableSchema.Schema
	}
	if tableSchema.Title != "" {
		setMap["title"] = tableSchema.Title
	}
	if tableSchema.FieldLen != 0 {
		setMap["field_len"] = tableSchema.FieldLen
	}
	if tableSchema.Description != "" {
		setMap["description"] = tableSchema.Description
	}
	if tableSchema.UpdatedAt != 0 {
		setMap["updated_at"] = tableSchema.UpdatedAt
	}
	if tableSchema.EditorID != "" {
		setMap["editor_id"] = tableSchema.EditorID
	}
	if tableSchema.EditorName != "" {
		setMap["editor_name"] = tableSchema.EditorName
	}
	return db.Table(t.TableName()).Where("app_id = ? and table_id = ? ", appID, tableID).Updates(
		setMap).Error
}
