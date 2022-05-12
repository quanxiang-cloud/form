package router

import (
	"net/http"

	"github.com/quanxiang-cloud/form/pkg/httputil"

	"github.com/labstack/echo/v4"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/form/internal/permit"
	echo2 "github.com/quanxiang-cloud/form/pkg/misc/echo"
)

func Permit(form permit.Permit) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &permit.Request{
			Request:  c.Request(),
			Response: c.Response(),
		}
		ctx := echo2.MutateContext(c)
		if err := bindParams(c, req); err != nil {
			logger.Logger.WithName("bind params").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)

			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		resp, err := form.Do(ctx, req)
		if err != nil {
			return err
		}

		if resp == nil {
			return echo.NewHTTPError(http.StatusForbidden)
		}

		return nil
	}
}

func bindParams(c echo.Context, i *permit.Request) error {
	if err := httputil.GetRequestArgs(c, &i.Data); err != nil {
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
