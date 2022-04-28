package form

import (
	"github.com/quanxiang-cloud/form/internal/service/consensus"
	"github.com/quanxiang-cloud/form/internal/service/types"
)

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
	Aggs  types.Any `json:"aggs"`
}

type CreateReq struct {
	Base
	Entity consensus.Entity `json:"entity"`
	Ref    types.Ref        `json:"ref"`
}

type GetReq struct {
	Base
	Query types.Query `json:"query"`
	Ref   types.Ref   `json:"ref"`
	Aggs  types.Any   `json:"aggs"`
}

type UpdateReq struct {
	Base
	Entity consensus.Entity `json:"entity"`
	Ref    types.Ref        `json:"ref"`
	Query  types.Query      `json:"query"`
}

type DeleteReq struct {
	Base
	Query types.Query `json:"query"`
}
