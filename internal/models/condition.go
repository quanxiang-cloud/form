package models

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/quanxiang-cloud/form/internal/service/types"
)

type Condition struct {
	Bool BOOL `json:"bool,omitempty"`
}

type BOOL struct {
	Must   types.Entities `json:"must,omitempty"`
	Should types.Entities `json:"should,omitempty"`
}

// Value 实现方法
func (c *Condition) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// Scan 实现方法
func (c *Condition) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &c)
}
