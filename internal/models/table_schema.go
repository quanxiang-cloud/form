package models

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
)

// SourceType SourceType.
type SourceType int64

const (
	FormSource SourceType = 1
)

// TableSchema TableSchema.
type TableSchema struct {
	ID          string
	AppID       string
	TableID     string
	FieldLen    int64
	Title       string
	Description string
	Source      SourceType
	CreatedAt   int64
	UpdatedAt   int64
	CreatorID   string
	CreatorName string
	EditorID    string
	EditorName  string
	Schema      SchemaProperties
}

type SchemaProperties map[string]SchemaProps

type SchemaProps struct {
	Title      string           `json:"title,omitempty"`
	Required   bool             `json:"required,omitempty"`
	Length     int              `json:"length,omitempty"`
	Type       string           `json:"type,omitempty"`
	ReadOnly   bool             `json:"read_only,omitempty"`
	Items      *SchemaProps     `json:"items,omitempty"`
	Properties SchemaProperties `json:"properties,omitempty"`
}

// Value 实现方法.
func (p SchemaProperties) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan 实现方法.
func (p *SchemaProperties) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &p)
}

type TableSchemaQuery struct {
	TableID    string
	AppID      string
	Title      string
	SourceType SourceType
}

type TableSchemeRepo interface {
	BatchCreate(db *gorm.DB, schema ...*TableSchema) error
	Get(db *gorm.DB, appID, tableID string) (*TableSchema, error)
	Delete(db *gorm.DB, query *TableSchemaQuery) error
	Update(db *gorm.DB, appID, tableID string, baseSchema *TableSchema) error
	List(db *gorm.DB, query *TableSchemaQuery, page, size int) ([]*TableSchema, int64, error)
}
