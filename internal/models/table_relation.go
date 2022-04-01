package models

import (
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
	Filter []string
}

type TableRelationQuery struct {
	AppID   string
	TableID string
}

type TableRelationRepo interface {
	BatchCreate(db *gorm.DB, table ...*TableRelation) error
	Find(db *gorm.DB, query *TableRelationQuery) ([]*TableRelation, error)
	Update(db *gorm.DB, table *TableRelation) error
	Delete(db *gorm.DB, query *TableRelationQuery) error
	List(db *gorm.DB, query *TableRelationQuery, page, size int) ([]*TableRelation, int64, error)
}
