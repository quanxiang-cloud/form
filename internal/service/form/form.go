package form

import (
	"context"

	"github.com/quanxiang-cloud/form/internal/service/types"
)

type Form interface {
	Search(ctx context.Context, req *SearchReq) (*SearchResp, error)
}

type SearchReq struct {
	AppID   string
	TableID string
	Page    int64
	Size    int64
	Sort    []string
	Query   types.Query
	Aggs    types.Any
	UserID  string
	DepID   string
}

type SearchResp struct {
	Aggregations interface{}              `json:"aggregations"`
	Entities     []map[string]interface{} `json:"entities"`
	Total        int64                    `json:"total"`
}
