package consensus

import (
	"context"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/internal/service/types"
)

type Universal struct {
	UserID string `json:"userID,omitempty"`
	// fixme
	DepID string `json:"depID,omitempty"`

	UserName string `json:"userName "`
}

type Foundation struct {
	TableID string `json:"tableID,omitempty"`
	AppID   string `json:"appID,omitempty"`

	Method string `json:"method,omitempty"`
}

type Get struct {
	Query types.Query `json:"query,omitempty"`
}

type List struct {
	Page int64     `json:"page,omitempty"`
	Size int64     `json:"size,omitempty"`
	Sort []string  `json:"sort,omitempty"`
	Aggs types.Any `json:"aggs,omitempty"`
}

type CreatedOrUpdate struct {
	Entity interface{} `json:"entity,omitempty"`
	Ref    types.Ref   `json:"ref"`
}

type Delete struct {
}

type Incidental struct {
	Permit *Permit `json:"-,omitempty"`
}

type Permit struct {
	Name      string
	Params    models.FiledPermit
	Response  models.FiledPermit
	Condition interface{}
	types     models.RoleType
}

type Bus struct {
	Universal
	Foundation
	Incidental
	Get
	List
	CreatedOrUpdate
	Delete
}
type Response interface{}

type Guidance interface {
	Do(ctx context.Context, bus *Bus) (Response, error)
}
