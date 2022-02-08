package models

import "context"

// OperationType OperationType
type OperationType int

const (
	// FindOperation FindOperation
	FindOperation OperationType = -1
	// DeleteInOperation DeleteInOperation
	DeleteInOperation OperationType = 1
	// DeleteEqualOperation DeleteEqualOperation
	DeleteEqualOperation OperationType = 2
	// DropOperation DropOperation
	DropOperation OperationType = 3
)

// Common  Common
type Common interface {
	ConstructSQL(ctx context.Context, table string, operation OperationType, column string, condition interface{}) string
	GetDBType(ctx context.Context) string
}
