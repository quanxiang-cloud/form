package models

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type SubTable struct {
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

type SubTaleQuery struct {
	AppID   string
	TableID string
}

type SubTableRepo interface {
	BatchCreate(ctx context.Context, db *mongo.Database, table ...*SubTable) error
	Find(ctx context.Context, db *mongo.Database, query *SubTaleQuery) ([]*SubTable, error)
	Update(ctx context.Context, db *mongo.Database, table *SubTable) error
	Delete(ctx context.Context, db *mongo.Database, query *SubTaleQuery) error
}
