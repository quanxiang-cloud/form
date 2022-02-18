package models

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
)

// SourceType SourceType
type SourceType int64

const (
	FormSource  SourceType = 1
	ModelSource SourceType = 2
)

// TableSchema TableSchema
type TableSchema struct {
	ID       string
	AppID    string
	TableID  string
	Title    string
	FieldLen int64

	Description string
	Source      SourceType
	CreatedAt   int64
	UpdatedAt   int64
	CreatorID   string
	CreatorName string
	EditorID    string
	EditorName  string
	Schema      TableSchemas
}

type TableSchemas map[string]interface{}

type TableSchemaQuery struct {
}

// Value 实现方法
func (p TableSchemas) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan 实现方法
func (p *TableSchemas) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &p)
}

type TableSchemeRepo interface {
	BatchCreate(db *gorm.DB, schema ...*TableSchema) error
	Get(db *gorm.DB, appID, tableID string) (*TableSchema, error)
	Find(db *gorm.DB, query *TableSchemaQuery, size int64, page int64) ([]*TableSchema, int64, error)
	Delete(db *gorm.DB, query *TableSchemaQuery) error
	Update(db *gorm.DB, appID, tableID string, baseSchema *TableSchema) error
}
