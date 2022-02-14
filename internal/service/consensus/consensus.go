package consensus

import (
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
	Ref   types.Ref   `json:"ref"`
}

type List struct {
	//Get  `json:"get,omitempty"`
	Page int64     `json:"page,omitempty"`
	Size int64     `json:"size,omitempty"`
	Sort []string  `json:"sort,omitempty"`
	Aggs types.Any `json:"aggs,omitempty"`
}

type CreatedOrUpdate struct {
	Entity interface{} `json:"entity,omitempty"`
}

type Incidental struct {
	Permit *Permit `json:"-,omitempty"`
}

type Permit struct {
	Filter      map[string]interface{}  `json:"filter,omitempty"`
	Condition   map[string]models.Query `json:"condition,omitempty"`
	Authority   int64                   `json:"authority,omitempty"`
	PermitTypes models.PerType          `json:"permitTypes,omitempty"`
}

type Bus struct {
	Universal
	Foundation
	Incidental
	Get
	List
	CreatedOrUpdate
}

type Response interface{}
