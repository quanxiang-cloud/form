package main

import (
	"database/sql/driver"
	"encoding/json"
)

// Table schema info
type Table struct {
	// id pk
	ID string `json:"id" bson:"_id"`

	AppID string `json:"appID" bson:"app_id"`
	// table id
	TableID string `json:"tableID" bson:"table_id"`
	// table design json schema
	Schema map[string]interface{} `json:"schema" bson:"schema"`
	// table page config json schema
	Config map[string]interface{} `json:"config" bson:"config"`
}

// SubTable table relation info
type SubTable struct {
	// id pk
	ID string `json:"id" bson:"_id"`
	// app id
	AppID string `json:"appID" bson:"app_id"`
	// table id
	TableID string `json:"tableID" bson:"table_id"`
	// table key name
	FieldName string `json:"fieldName" bson:"field_name"`
	// sub table id
	SubTableID string `json:"subTableID" bson:"sub_table_id"`
	// table type
	SubTableType string `json:"subTableType" bson:"sub_table_type"`
	// filter
	Filter []string `json:"filter" bson:"filter"`
}

// DataBaseSchema DataBaseSchema
type DataBaseSchema struct {
	ID          string                 `bson:"_id"`
	Title       string                 `bson:"title"`
	AppID       string                 `bson:"app_id"`
	TableID     string                 `bson:"table_id"`
	FieldLen    int64                  `bson:"field_len"`
	Description string                 `bson:"description"`
	Source      int                    `bson:"source"`
	CreatedAt   int64                  `bson:"created_at"`
	UpdatedAt   int64                  `bson:"updated_at"`
	CreatorID   string                 `bson:"creator_id"`
	CreatorName string                 `bson:"creator_name"`
	EditorID    string                 `bson:"editor_id"`
	EditorName  string                 `bson:"editor_name"`
	Schema      map[string]interface{} `bson:"schema"` // 过滤之后的代码
}

// TableSchema TableSchema.
type TableSchema struct {
	ID          string
	AppID       string
	TableID     string
	FieldLen    int64
	Title       string
	Description string
	Source      int
	CreatedAt   int64
	UpdatedAt   int64
	CreatorID   string
	CreatorName string
	EditorID    string
	EditorName  string
	Schema      SchemaProperties
}

type SchemaProperties map[string]interface{}

// Value 实现方法.
func (p SchemaProperties) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan 实现方法.
func (p *SchemaProperties) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &p)
}
