package models

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
)

// Permit Permit
type Permit struct {
	ID        string
	Path      string
	Params    FiledPermit
	Response  FiledPermit
	Condition interface{}

	RoleID string

	CreatedAt   int64
	CreatorID   string
	CreatorName string
}

type Bool map[string][]Query

type Aggs map[string]Agg

type Agg map[string]struct {
	Field string `json:"field"`
}

type Query map[string]Field

type Field map[string]Value

type Value interface{}

type DSL struct {
	QY map[string]interface{} `json:"query"`

	Query Query
	Bool  Bool

	Aggs Aggs `json:"aggs"`
}

//

type Key struct {
	Type       string
	Properties FiledPermit
}

type FiledPermit map[string]Key

// Value 实现方法
func (p FiledPermit) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan 实现方法
func (p *FiledPermit) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &p)
}

type PermitQuery struct {
	RoleID string
}

// PermitRepo PermitRepo
type PermitRepo interface {
	BatchCreate(db *gorm.DB, form ...*Permit) error

	Get(db *gorm.DB, path, roleID string) (*Permit, error)

	Find(db *gorm.DB, query *PermitQuery) ([]*Permit, error)

	Delete(db *gorm.DB, query *PermitQuery) error

	Update(db *gorm.DB, path, roleID string, permit *Permit) error
}
