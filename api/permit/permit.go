package router

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/form/internal/permit"
)

func ProxyForm(form permit.Permit) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &permit.Request{
			Request: c.Request(),
			Writer:  c.Response().Writer,
		}
		if err := bindParams(c, req); err != nil {
			logger.Logger.Errorw("bind request body param error", "error", err)
			c.NoContent(http.StatusBadRequest)
			return nil
		}

		resp, err := form.Do(context.Background(), req)
		if err != nil {
			return err
		}
		if resp == nil {
			c.NoContent(http.StatusForbidden)
			return nil
		}

		fmt.Println(resp)

		return nil
	}
}

func ProxyPoly(poly permit.Permit) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &permit.Request{
			Request: c.Request(),
			Writer:  c.Response().Writer,
		}

		if err := bindParams(c, req); err != nil {
			logger.Logger.Errorw("bind request body param error", "error", err)
			c.JSON(http.StatusBadRequest, err)
			return nil
		}

		resp, err := poly.Do(context.Background(), req)
		if err != nil {
			return err
		}

		if resp == nil {
			c.NoContent(http.StatusForbidden)
		}

		return nil
	}
}

func bindParams(c echo.Context, i *permit.Request) error {
	if err := (&echo.DefaultBinder{}).BindBody(c, &i.Body); err != nil {
		logger.Logger.Errorw("bind request body param error", "error", err)
		return err
	}

	if err := (&echo.DefaultBinder{}).BindQueryParams(c, i); err != nil {
		logger.Logger.Errorw("bind request query param error", "error", err)
		return err
	}

	if err := (&echo.DefaultBinder{}).BindPathParams(c, i); err != nil {
		logger.Logger.Errorw("bind request path param error", "error", err)
		return err
	}

	if err := (&echo.DefaultBinder{}).BindHeaders(c, i); err != nil {
		logger.Logger.Errorw("bind request header param error", "error", err)
		return err
	}

	return nil
}
