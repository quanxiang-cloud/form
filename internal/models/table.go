package models

import (
	"database/sql/driver"
	"encoding/json"

	"gorm.io/gorm"
)

// Table schema info
type Table struct {
	// id pk
	ID string

	AppID string
	// table id
	TableID string

	// table design json schema
	Schema WebSchema
	// table page config json schema
	Config Config
}

//WebSchema WebSchema
type WebSchema map[string]interface{}

// Config Config
type Config map[string]interface{}

// Value 实现方法
func (p WebSchema) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan 实现方法
func (p *WebSchema) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &p)
}

// Value 实现方法
func (p Config) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan 实现方法
func (p *Config) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &p)
}

type TableQuery struct {
	AppID    string
	TableID  string
	TableIDS []string
}
type TableRepo interface {
	BatchCreate(db *gorm.DB, tables ...*Table) error
	Get(db *gorm.DB, appId, tableID string) (*Table, error)
	Find(db *gorm.DB, query *TableQuery) ([]*Table, error)
	Delete(db *gorm.DB, query *TableQuery) error
	Update(db *gorm.DB, appID, tableID string, table *Table) error
}
