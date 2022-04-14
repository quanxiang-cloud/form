package router

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/form/internal/component/event"
	"github.com/quanxiang-cloud/form/internal/permit"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
	echo2 "github.com/quanxiang-cloud/form/pkg/misc/echo"
)

// UserRole user role.
func (p *Cache) UserRole(c echo.Context) error {
	ctx := echo2.MutateContext(c)

	data := new(event.DaprEvent)
	if err := c.Bind(data); err != nil {
		logger.Logger.WithName("bind params").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.NoContent(http.StatusBadRequest)

		return nil
	}

	if data.Data.UserSpec == nil {
		data.Data.UserSpec = &event.UserSpec{}
	}

	req := &permit.UserMatchReq{
		UserID: data.Data.UserSpec.UserID,
		RoleID: data.Data.UserSpec.RoleID,
		AppID:  data.Data.UserSpec.AppID,
		Action: data.Data.UserSpec.Action,
	}
	_, err := p.cache.UserRole(context.Background(), req)
	if err != nil {
		logger.Logger.WithName("user role").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		return err
	}

	return nil
}

// Permit permit.
func (p *Cache) Permit(c echo.Context) error {
	ctx := echo2.MutateContext(c)

	data := new(event.DaprEvent)
	if err := c.Bind(data); err != nil {
		logger.Logger.WithName("bind params").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		c.NoContent(http.StatusBadRequest)

		return nil
	}

	per := data.Data.PermitSpec
	req := &permit.LimitReq{
		RoleID:    per.RoleID,
		Path:      per.Path,
		Condition: per.Condition,
		Params:    per.Params,
		Response:  per.Response,
		Action:    per.Action,
	}
	_, err := p.cache.Limit(context.Background(), req)
	if err != nil {
		logger.Logger.WithName("user role").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		return err
	}

	return nil
}

// Cache cache.
type Cache struct {
	cache permit.Cache
}

// NewCache new cache.
func NewCache(config *config.Config) (*Cache, error) {
	newCache, err := permit.NewCache(config)
	if err != nil {
		return nil, err
	}

	return &Cache{
		cache: newCache,
	}, nil
}
