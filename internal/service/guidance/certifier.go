package guidance

import (
	"context"

	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/form/internal/filters"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/internal/service"
	"github.com/quanxiang-cloud/form/internal/service/consensus"
	"github.com/quanxiang-cloud/form/pkg/misc/code"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

type certifier struct {
	permit service.Permission

	next Guidance
}

func newCertifier(conf *config.Config) (Guidance, error) {
	permit, err := service.NewPermission(conf)
	if err != nil {
		return nil, err
	}

	next, err := newRuling()
	if err != nil {
		return nil, err
	}
	return &certifier{
		permit: permit,

		next: next,
	}, nil
}

func (c *certifier) Do(ctx context.Context, bus *consensus.Bus) (consensus.Response, error) {
	// err := c.pre(ctx, bus, checkOperate())
	// if err != nil {
	// 	return nil, err
	// }

	resp, err := c.next.Do(ctx, bus)
	if err != nil {
		return nil, err
	}

	// FIXME 不应该关心返回结构
	// c.post(ctx, bus, JSONFilter())
	return resp, err
}

func (c *certifier) pre(ctx context.Context, bus *consensus.Bus, opts ...preOption) error {
	// get permit
	err := c.getPermit(ctx, bus)
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

func (c *certifier) getPermit(ctx context.Context, bus *consensus.Bus) error {
	cache, err := c.permit.GetPerInCache(ctx, &service.GetPerInCacheReq{
		UserID: bus.UserID,
		DepID:  bus.DepID,
		FormID: bus.TableID,
		AppID:  bus.AppID,
	})
	if err != nil {
		return err
	}
	if cache == nil {
		return error2.New(code.ErrNotPermit)
	}

	bus.Permit.Condition = cache.DataAccessPer
	bus.Permit.Authority = cache.Authority
	bus.Permit.Filter = cache.Filter
	bus.Permit.PermitTypes = cache.Type

	return nil
}

type preOption func(ctx context.Context, bus *consensus.Bus) bool

func checkOperate() preOption {
	return func(ctx context.Context, bus *consensus.Bus) bool {
		if bus.Permit.PermitTypes == models.InitType {
			return true
		}
		var op int64
		switch bus.Method {
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

		return op&bus.Permit.Authority != 0
	}
}

func CheckData() preOption {
	return func(ctx context.Context, bus *consensus.Bus) bool {
		if bus.Entity == nil {
			return false
		}
		if bus.Permit.PermitTypes == models.InitType {
			return true
		}
		return filters.FilterCheckData(bus.Entity, bus.Permit.Filter)
	}
}

//FilterOption FilterOption
type FilterOption func(data interface{}, filter map[string]interface{}) error

func (c *certifier) post(ctx context.Context, bus *consensus.Bus, opts ...FilterOption) error {
	if bus.Permit.PermitTypes == models.InitType {
		return nil
	}
	for _, opt := range opts {
		opt(bus.Entity, bus.Permit.Filter)
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
