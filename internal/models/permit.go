package models

import (
	"database/sql/driver"
	"encoding/json"

	"gorm.io/gorm"
)

// Permit Permit
type Permit struct {
	ID          string
	RoleID      string
	Path        string
	Params      FiledPermit
	Response    FiledPermit
	Condition   Condition
	CreatedAt   int64
	CreatorID   string
	CreatorName string
}

//

type Key struct {
	Type       string      `json:"type,omitempty"`
	Properties FiledPermit `json:"properties,omitempty"`
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
	ID     string
	RoleID string
	Path   string
}

// PermitRepo PermitRepo
type PermitRepo interface {
	BatchCreate(db *gorm.DB, form ...*Permit) error

	Get(db *gorm.DB, roleID, path string) (*Permit, error)

	Find(db *gorm.DB, query *PermitQuery) ([]*Permit, error)

	Delete(db *gorm.DB, query *PermitQuery) error

	Update(db *gorm.DB, id string, permit *Permit) error
}
