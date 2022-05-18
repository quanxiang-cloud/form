package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/form/internal/permit"
	"github.com/quanxiang-cloud/form/pkg/httputil"
	echo2 "github.com/quanxiang-cloud/form/pkg/misc/echo"

)

const (
	path = "path"
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
	i.Path = c.Request().URL.String()
	if c.Get(path) != "" {
		s, ok := c.Get(path).(string)
		if ok {
			i.Path = s
		}
	}
	return nil
}

// V2FormPath ,V2FormPath
func V2FormPath(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		paths := c.Request().URL.String()
		paths = fmt.Sprintf("%s/:id", paths[0:strings.LastIndex(paths, "/")])
		c.Set(path, paths)
		fmt.Println(paths)
		return next(c)
	}
}
