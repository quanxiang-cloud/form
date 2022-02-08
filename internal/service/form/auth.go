package form

import (
	"context"
	config2 "github.com/quanxiang-cloud/form/pkg/misc/config"

	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/form/internal/filters"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/internal/service"
	"github.com/quanxiang-cloud/form/pkg/misc/code"
)

type auth struct {
	permit service.Permission
	*comet
}

func NewAuthForm(conf *config2.Config) (Form, error) {
	permits, err := service.NewPermission(conf)
	if err != nil {
		return nil, err
	}
	comet, err := newForm()
	if err != nil {
		return nil, err
	}
	return &auth{
		permit: permits,
		comet:  comet,
	}, nil
}

func (a *auth) Search(ctx context.Context, req *SearchReq) (*SearchResp, error) {

	bus := &bus{
		AppID:     req.AppID,
		userID:    req.UserID,
		depID:     req.DepID,
		tableName: req.TableID,
		permit:    &permit{},
		method:    "find",
	}
	err := a.pre(ctx, bus, checkOperate())
	if err != nil {
		return nil, err
	}
	resp, err := a.comet.Search(ctx, req)
	if err != nil {
		return nil, err
	}
	bus.entity = resp.Entities
	a.post(ctx, bus, JSONFilter())
	return resp, nil
}

type bus struct {
	userID    string
	depID     string
	tableName string
	AppID     string
	permit    *permit
	entity    interface{}
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

//FilterOption FilterOption
type FilterOption func(data interface{}, filter map[string]interface{}) error

func (a *auth) post(ctx context.Context, bus *bus, opts ...FilterOption) error {
	if bus.permit.permitTypes == models.InitType {
		return nil
	}
	for _, opt := range opts {
		opt(bus.entity, bus.permit.Filter)
	}
	return nil
}

// JSONFilter JSONFilter
func JSONFilter() FilterOption {
	return func(data interface{}, filter map[string]interface{}) error {
		if data == nil {
			return nil
		}

		if filter == nil {
			return error2.New(code.ErrNotPermit)
		}
		filters.JSONFilter2(data, filter)
		return nil
	}
}
