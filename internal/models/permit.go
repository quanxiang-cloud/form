package models

import (
	"database/sql/driver"
	"encoding/json"

	"gorm.io/gorm"
)

// Permit Permit.
type Permit struct {
	ID          string
	RoleID      string
	Path        string
	Params      FiledPermit
	Response    FiledPermit
	Condition   Condition
	Method      string
	ParamsAll   bool
	ResponseAll bool
	CreatedAt   int64
	CreatorID   string
	CreatorName string
}

type Condition map[string]interface{}

// Value 实现方法
func (c Condition) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// Scan 实现方法
func (c *Condition) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &c)
}

type Key struct {
	Type       string      `json:"type,omitempty"`
	Properties FiledPermit `json:"properties,omitempty"`
}

type FiledPermit map[string]Key

// Value 实现方法.
func (p FiledPermit) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan 实现方法.
func (p *FiledPermit) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &p)
}

type PermitQuery struct {
	ID      string
	RoleID  string
	Path    string
	Method  string
	Paths   []string
	RoleIDs []string
}

// PermitRepo PermitRepo.
type PermitRepo interface {
	BatchCreate(db *gorm.DB, form ...*Permit) error

	Get(db *gorm.DB, roleID, path, method string) (*Permit, error)

	Delete(db *gorm.DB, query *PermitQuery) error

	Update(db *gorm.DB, query *PermitQuery, permit *Permit) error

	List(db *gorm.DB, query *PermitQuery, page, size int) ([]*Permit, int64, error)
}
