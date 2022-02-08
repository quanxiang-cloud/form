package form

import (
	"context"

	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/form/internal/filters"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/internal/service"
	"github.com/quanxiang-cloud/form/pkg/misc/code"
)

type auth struct {
	permit service.Permission
	comet
}

func NewAuthForm() Form {
	return &auth{}
}

func (a *auth) Search(ctx context.Context, req *SearchReq) (*SearchResp, error) {
	bus := &bus{
		// TODO
	}
	a.pre(ctx, bus, checkOperate())
	resp, err := a.comet.Search(ctx, req)
	if err != nil {
		return nil, err
	}
	a.post(ctx, bus)
	return resp, nil
}

type bus struct {
	userID    string
	depID     string
	tableName string
	AppID     string
	permit    *permit
	entity    *Entity
	method    string
}

type permit struct {
	Filter      map[string]interface{}
	Condition   map[string]models.Query
	Authority   int64
	permitTypes models.PerType
}

type preOption func(ctx context.Context, bus *bus) bool

func checkOperate() preOption {
	return func(ctx context.Context, bus *bus) bool {
		if bus.permit.permitTypes == models.InitType {
			return true
		}
		var op int64
		switch bus.method {
		case "find", "findOne":
			op = models.OPRead
		case "create":
			op = models.OPCreate
		case "update":
			op = models.OPUpdate
		case "delete":
			op = models.OPDelete
		default:
			return false
		}

		return op&bus.permit.Authority != 0
	}
}

func CheckData() preOption {
	return func(ctx context.Context, bus *bus) bool {
		if bus.entity == nil {
			return false
		}
		if bus.permit.permitTypes == models.InitType {
			return true
		}
		return filters.FilterCheckData(bus.entity, bus.permit.Filter)
	}
}

func (a *auth) pre(ctx context.Context, bus *bus, opts ...preOption) error {
	// get permit
	err := a.getPermit(ctx, bus)
	if err != nil {
		return err
	}

	for _, opt := range opts {
		if !opt(ctx, bus) {
			return error2.New(code.ErrNotPermit)
		}
	}

	return nil
}

func (a *auth) getPermit(ctx context.Context, bus *bus) error {
	cache, err := a.permit.GetPerInCache(ctx, &service.GetPerInCacheReq{
		UserID: bus.userID,
		DepID:  bus.depID,
		FormID: bus.tableName,
		AppID:  bus.AppID,
	})
	if err != nil {
		return err
	}
	if cache == nil {
		return error2.New(code.ErrNotPermit)
	}

	bus.permit.Condition = cache.DataAccessPer
	bus.permit.Authority = cache.Authority
	bus.permit.Filter = cache.Filter
	bus.permit.permitTypes = cache.Type

	return nil
}

func (a *auth) post(ctx context.Context, bus *bus) {}
