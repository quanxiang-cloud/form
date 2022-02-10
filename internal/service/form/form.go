package form

import (
	"context"

	"github.com/quanxiang-cloud/form/internal/service/types"
)

type Form interface {
	Search(ctx context.Context, req *SearchReq) (*SearchResp, error)
	Create(ctx context.Context, req *CreateReq) (*CreateResp, error)
}

type base struct {
	AppID   string
	TableID string

	UserID string
	DepID  string
}

type SearchReq struct {
	base
	Page  int64
	Size  int64
	Sort  []string
	Query types.Query
	Aggs  types.Any
}

type SearchResp struct {
	Aggregations interface{}              `json:"aggregations"`
	Entities     []map[string]interface{} `json:"entities"`
	Total        int64                    `json:"total"`
}

type CreateReq struct {
	AppID       string    `json:"appID"`
	TableID     string    `json:"tableID"`
	Entity      Entity    `json:"entity"`
	Ref         types.Ref `json:"ref"`
	UserID      string    `json:"userID"`
	DepID       string    `json:"depID"`
	CreatorName string    `json:"creatorName"`
}

type CreateResp struct {
	Entity     interface{} `json:"entity"`
	ErrorCount int64       `json:"errorCount"`
}
