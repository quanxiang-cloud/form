package router

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/quanxiang-cloud/form/pkg/misc/config"

	"github.com/labstack/echo/v4"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/form/internal/component/event"
	"github.com/quanxiang-cloud/form/internal/permit"
)

func (p *Cache) UserRole(c echo.Context) error {
	data := new(event.DaprEvent)
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(body, &data)
	if err != nil {

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
	_, err = p.cache.UserRole(context.Background(), req)
	if err != nil {
		logger.Logger.Errorw("msg is error ", err.Error())
		return err
	}
	return nil
}

func (p *Cache) Permit(c echo.Context) error {
	data := new(event.DaprEvent)
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(body, &data)
	if err != nil {

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
	_, err = p.cache.Limit(context.Background(), req)
	if err != nil {
		logger.Logger.Errorw("msg is error ", err.Error())
		return err
	}
	return nil
}

type Cache struct {
	cache permit.Cache
}

func NewCache(config *config.Config) (*Cache, error) {
	newCache, err := permit.NewCache(config)
	if err != nil {
		return nil, err
	}
	return &Cache{
		cache: newCache,
	}, nil
}
