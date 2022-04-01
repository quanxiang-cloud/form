package models

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
)

type TableRelation struct {
	ID string
	// app id
	AppID string
	// table id
	TableID string
	// table key name
	FieldName string
	// sub table id
	SubTableID string
	// table type
	SubTableType string
	// filter
	Filter Filters
}

type Filters []string

// Value 实现方法
func (p Filters) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan 实现方法
func (p *Filters) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &p)
}

type TableRelationQuery struct {
	AppID   string
	TableID string
}

type TableRelationRepo interface {
	BatchCreate(db *gorm.DB, table ...*TableRelation) error
	Find(db *gorm.DB, query *TableRelationQuery) ([]*TableRelation, error)
	Update(db *gorm.DB, tableID, fieldName string, table *TableRelation) error
	Delete(db *gorm.DB, query *TableRelationQuery) error
	Get(db *gorm.DB, tableID, fieldName string) (*TableRelation, error)
}
