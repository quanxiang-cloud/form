package models

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
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
	Insert(ctx context.Context, db *mongo.Database, dataset *DataSet) error
	Update(ctx context.Context, db *mongo.Database, dataset *DataSet) error
	Delete(ctx context.Context, db *mongo.Database, id string) error
	GetByID(ctx context.Context, db *mongo.Database, id string) (*DataSet, error)
	Find(ctx context.Context, db *mongo.Database, query *DataSetQuery) ([]*DataSet, error)
}
