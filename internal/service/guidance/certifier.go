package guidance

import (
	"context"
	"github.com/quanxiang-cloud/form/internal/filters"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/internal/service"

	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/form/internal/service/consensus"
	"github.com/quanxiang-cloud/form/pkg/misc/code"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

type certifier struct {
	permit service.Permit

	next consensus.Guidance
}

func newCertifier(conf *config.Config) (consensus.Guidance, error) {
	permit, err := service.NewPermit(conf)
	if err != nil {
		return nil, err
	}
	next, err := newRuling(conf)
	if err != nil {
		return nil, err
	}
	return &certifier{
		permit: permit,

		next: next,
	}, nil
}

func (c *certifier) Do(ctx context.Context, bus *consensus.Bus) (*consensus.Response, error) {
	//err := c.pre(ctx, bus)
	//if err != nil {
	//	return nil ,err
	//}
	resp, err := c.next.Do(ctx, bus)
	if err != nil {
		return nil, err
	}
	//c.post(ctx ,bus ,resp )
	return resp, err
}

func (c *certifier) pre(ctx context.Context, bus *consensus.Bus) error {
	// get permit
	err := c.getPermit(ctx, bus)
	if err != nil {
		return err
	}
	// 判断有无权限
	if bus.Permit.Types == models.InitType {
		return nil
	}

	if !filters.Pre(bus.Entity, bus.Permit.Params) {
		return error2.New(code.ErrNotPermit)
	}
	return nil
}

func (c *certifier) getPermit(ctx context.Context, bus *consensus.Bus) error {
	cache, err := c.permit.GetPerInCache(ctx, &service.GetPerInCacheReq{
		UserID: bus.UserID,
		DepID:  bus.DepID,
		Path:   bus.Path,
		AppID:  bus.AppID,
	})
	if err != nil {
		return err
	}
	if cache == nil {
		return error2.New(code.ErrNotPermit)
	}
	permits := &consensus.Permit{
		Params:    cache.Params,
		Response:  cache.Response,
		Condition: cache.Condition,
		Types:     cache.Types,
	}
	bus.Permit = permits
	return nil
}

func (c *certifier) post(ctx context.Context, bus *consensus.Bus, resp *consensus.Response) {
	var entity interface{}
	switch bus.Method {
	case "get":
		entity = resp.GetResp.Entity
	case "search":
		entity = resp.ListResp.Entities
	}
	filters.Post(entity, bus.Permit.Response)
}
