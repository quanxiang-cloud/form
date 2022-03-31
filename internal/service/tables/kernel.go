package tables

import (
	"context"

	"github.com/quanxiang-cloud/form/internal/models"
)

type Bus struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`

	AppID   string           `json:"app_id"`
	TableID string           `json:"tableID"`
	Schema  models.WebSchema `json:"schema"`

	Source models.SourceType `json:"source"` // source 1 是表单驱动，2是模型驱动
	Update bool              `json:"update"`
	ConvertSchemas
}

type ConvertSchemas struct {
	ConvertSchema models.SchemaProperties `json:"convertSchema"`
	Title         string
	FieldLen      int64
	Description   string
}

type DoResponse struct {
}

type Guidance interface {
	Do(ctx context.Context, bus *Bus) (*DoResponse, error)
}
