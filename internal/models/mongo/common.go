package mongo

import (
	"context"
	"fmt"
	"github.com/quanxiang-cloud/form/internal/models"
)

const (
	collectionOperateTemp = "db.getCollection(\"%s\")"
	dropOperateTemp       = ".drop()"
	deleteInFilterTemp    = ".deleteMany({\"%s\":{\"$in\":%v}})"
	deleteEqualFilterTemp = ".deleteMany({\"%s\": \"%s\"})"
	findFilterTemp        = ".find({\"%s\":\"%s\"})"
	dbType                = "mongoDB"
)

type common struct {
}

// NewCommon NewCommon
func NewCommon() models.Common {
	return &common{}
}

func (c *common) ConstructSQL(ctx context.Context, table string, operation models.OperationType, column string, condition interface{}) string {
	sql := fmt.Sprintf(collectionOperateTemp, table)
	switch operation {
	case models.DeleteInOperation:
		sql += fmt.Sprintf(deleteInFilterTemp, column, condition)
	case models.DeleteEqualOperation:
		sql += fmt.Sprintf(deleteEqualFilterTemp, column, condition)
	case models.DropOperation:
		sql += dropOperateTemp
	case models.FindOperation:
		sql += fmt.Sprintf(findFilterTemp, column, condition)
	}
	return sql
}

func (c *common) GetDBType(ctx context.Context) string {
	return dbType
}
