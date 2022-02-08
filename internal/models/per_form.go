package models

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
)

const (
	// OPRead find or findOne
	OPRead = 1 << iota
	// OPCreate create
	OPCreate
	// OPUpdate update
	OPUpdate
	// OPDelete delete
	OPDelete
	// OPImport import
	OPImport
	// OPExport export
	OPExport
)

// PermitForm PermitForm
type PermitForm struct {
	PerGroupID string `gorm:"column:permit_id"`
	FormID     string
	FormType   string
	Authority  int64
	Conditions Conditions
	FieldJSON  FieldJSON
	WebSchema  Schema
}

// Schema schema
type Schema struct {
	Title      string            `json:"title,omitempty"`
	Types      string            `json:"type,omitempty"`
	XInternal  XInternal         `json:"x-internal,omitempty"`
	Properties map[string]Schema `json:"properties,omitempty"` // type==object
	Item       *Schema           `json:"item,omitempty"`       //type==array
}

// XInternal x-internal
type XInternal struct {
	Sortable   bool    `json:"sortable"`   //排序
	Permission float64 `json:"permission"` //权限属性，第一位可不可见，第二位可不可编辑
}

type PerFormQuery struct {
	PerGroupID string
	FormID     string
	GroupIDS   string
}

type Conditions map[string]Query

// Query Query
type Query map[string]interface{}

// FieldJSON FieldJSON
type FieldJSON map[string]interface{}

// Value 实现方法
func (p Conditions) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan 实现方法
func (p *Conditions) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &p)
}

// Scan 实现方法
func (p *Schema) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &p)
}

// Value 实现方法
func (p Schema) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Value 实现方法
func (p FieldJSON) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan 实现方法
func (p *FieldJSON) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &p)
}

// GroupFormRepo GroupFormRepo
type GroupFormRepo interface {
	BatchCreate(db *gorm.DB, form ...*PermitForm) error
	Get(db *gorm.DB, permitID, formID string) (*PermitForm, error)
	Find(db *gorm.DB, query *PerFormQuery) ([]*PermitForm, error)
	Delete(db *gorm.DB, query *PerFormQuery) error
	Update(db *gorm.DB, permitID, formID string, form *PermitForm) error
}
