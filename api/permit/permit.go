package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/form/internal/permit"
	echo2 "github.com/quanxiang-cloud/form/pkg/misc/echo"
)

// ProxyForm form proxy.
func ProxyForm(form permit.Permit) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &permit.Request{
			Request: c.Request(),
			Writer:  c.Response().Writer,
		}

		ctx := echo2.MutateContext(c)
		if err := bindParams(c, req); err != nil {
			logger.Logger.WithName("bind params").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
			c.NoContent(http.StatusBadRequest)

			return nil
		}

		resp, err := form.Do(ctx, req)
		if err != nil {
			return err
		}

		if resp == nil {
			c.NoContent(http.StatusForbidden)
			return nil
		}

		return nil
	}
}

// ProxyPoly polyapi proxy.
func ProxyPoly(poly permit.Permit) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &permit.Request{
			Request: c.Request(),
			Writer:  c.Response().Writer,
		}

		ctx := echo2.MutateContext(c)
		if err := bindParams(c, req); err != nil {
			logger.Logger.WithName("bind params").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
			c.JSON(http.StatusBadRequest, err)
			return nil
		}

		resp, err := poly.Do(ctx, req)
		if err != nil {
			return err
		}

		if resp == nil {
			c.NoContent(http.StatusForbidden)
			return nil
		}

		return nil
	}
}

func bindParams(c echo.Context, i *permit.Request) error {
	var err error
	switch c.Request().Method {
	case http.MethodPost:
		i.Body = make(permit.Object)
		err = (&echo.DefaultBinder{}).BindBody(c, &i.Body)
	case http.MethodGet:
		i.Query, i.Entity = make(permit.Object), make(permit.Object)
		err = (&echo.DefaultBinder{}).BindQueryParams(c, i)
	}

	if err != nil {
		return err
	}

	if err := (&echo.DefaultBinder{}).BindPathParams(c, i); err != nil {
		return err
	}

	if err := (&echo.DefaultBinder{}).BindHeaders(c, i); err != nil {
		return err
	}

	return nil
}
