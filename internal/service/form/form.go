package form

import (
	"context"

	"github.com/quanxiang-cloud/form/internal/service/types"
)

type Form interface {
	Search(ctx context.Context, req *SearchReq) (*SearchResp, error)
	Create(ctx context.Context, req *CreateReq) (*CreateResp, error)
	Get(ctx context.Context, req *GetReq) (*GetResp, error)
	Update(ctx context.Context, req *UpdateReq) (*UpdateResp, error)
	Delete(ctx context.Context, req *DeleteReq) (*DeleteResp, error)
}

type Base struct {
	AppID   string
	TableID string

	UserID   string
	DepID    string
	UserName string
}

type SearchReq struct {
	Base
	Page  int64
	Size  int64
	Sort  []string
	Query types.Query
	Aggs  types.Any
}

type SearchResp struct {
	Entities types.Entities `json:"entities"`
	Total    int64          `json:"total"`
}

type CreateReq struct {
	Base
	Entity Entity    `json:"entity"`
	Ref    types.Ref `json:"ref"`
}

type CreateResp struct {
	Entity Entity `json:"entity"`
	Count  int64  `json:"errorCount"`
}

type GetReq struct {
	Base
	AppID string
	Query types.Query `json:"query"`
	Ref   types.Ref   `json:"ref"`
}

type GetResp struct {
	Entity types.M `json:"entity"`
}

type UpdateReq struct {
	Base
	Entity Entity      `json:"entity"`
	Ref    types.Ref   `json:"ref"`
	Query  types.Query `json:"query"`
}

type UpdateResp struct {
	Count int64 `json:"errorCount"`
}

type DeleteReq struct {
	Base
	Query types.Query `json:"query"`
}

type DeleteResp struct {
	Count int64 `json:"errorCount"`
}
