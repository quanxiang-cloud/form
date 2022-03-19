package models

import (
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
	ID      string
	AppID   string
	TableID string

	FieldLen int64
	Title    string

	Description string
	Source      SourceType
	CreatedAt   int64
	UpdatedAt   int64
	CreatorID   string
	CreatorName string
	EditorID    string
	EditorName  string
	Schema      WebSchema
}

type TableSchemaQuery struct {
	TableID string
	AppID   string
}

type TableSchemeRepo interface {
	BatchCreate(db *gorm.DB, schema ...*TableSchema) error
	Get(db *gorm.DB, appID, tableID string) (*TableSchema, error)
	Find(db *gorm.DB, query *TableSchemaQuery, size int64, page int64) ([]*TableSchema, int64, error)
	Delete(db *gorm.DB, query *TableSchemaQuery) error
	Update(db *gorm.DB, appID, tableID string, baseSchema *TableSchema) error
}
