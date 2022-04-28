package consensus

import (
	"context"

	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/internal/service/types"
)

type Universal struct {
	UserID string `json:"userID,omitempty"`

	DepID string `json:"depID,omitempty"`

	UserName string `json:"userName"`

	Path string `json:"path"`
}

type Ref struct {
	Ref types.Ref `json:"ref,omitempty"`
}

type Foundation struct {
	TableID string `json:"tableID,omitempty"`
	AppID   string `json:"appID,omitempty"`
	Method  string `json:"method,omitempty"`
}

type Get struct {
	Query    types.Query `json:"query,omitempty" form:"query"`
	OldQuery types.Query `json:"OldQuery"`
	Aggs     types.Any   `json:"aggs"`
}

type List struct {
	Page int64    `json:"page,omitempty" form:"page"`
	Size int64    `json:"size,omitempty" form:"size"`
	Sort []string `json:"sort,omitempty" form:"sort"`
}

type CreatedOrUpdate struct {
	Entity interface{} `json:"entity,omitempty"`
}

type Delete struct{}

type Incidental struct {
	Permit *Permit `json:"-,omitempty"`
}

type Permit struct {
	Name        string
	Params      models.FiledPermit
	Response    models.FiledPermit
	Condition   models.Condition
	Types       models.RoleType
	ResponseAll bool
	ParamsAll   bool
}

type Bus struct {
	Universal
	Foundation
	Incidental
	Ref
	Get
	List
	CreatedOrUpdate
	Delete
	Sub
}
type Sub struct {
	PID      string `json:"pid" form:"pid"`
	PTableID string `json:"pTableID" form:"pTableID"`
	FieldKey string `json:"fieldKey" form:"fieldKey"`
}

type Response struct {
	Entity   types.M        `json:"entity,omitempty"`
	Total    int64          `json:"total"`
	Entities types.Entities `json:"entities,omitempty"`
}
type Guidance interface {
	Do(ctx context.Context, bus *Bus) (*Response, error)
}
