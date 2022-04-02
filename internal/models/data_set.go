package models

import (
	"gorm.io/gorm"
)

// DataSet DataSet
type DataSet struct {
	ID        string `bson:"_id"`
	Name      string `bson:"name"`
	Tag       string `bson:"tag"`
	Type      int64  `bson:"type"`
	Content   string `bson:"content"`
	CreatedAt int64  `bson:"created_at"`
}

type DataSetQuery struct {
	Tag   string
	Name  string
	Types int64
}

// DataSetRepo 数据层接口
type DataSetRepo interface {
	Insert(db *gorm.DB, dataset *DataSet) error
	Update(db *gorm.DB, id string, dataset *DataSet) error
	Delete(db *gorm.DB, id string) error
	GetByID(db *gorm.DB, id string) (*DataSet, error)
	Find(db *gorm.DB, query *DataSetQuery, page, size int) ([]*DataSet, int64, error)
}
