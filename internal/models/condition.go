package models

import (
	"database/sql/driver"
	"encoding/json"
)

type Condition map[string]interface{}

// Value 实现方法
func (c *Condition) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// Scan 实现方法
func (c *Condition) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &c)
}
