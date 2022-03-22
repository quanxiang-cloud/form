package router

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/form/internal/permit"
)

func ProxyForm(form permit.Form) echo.HandlerFunc {
	return func(c echo.Context) error {
		guardReq := &permit.GuardReq{
			Request: c.Request(),
			Writer:  c.Response().Writer,
			Body:    map[string]interface{}{},
		}
		if err := bindParams(c, guardReq); err != nil {
			logger.Logger.Errorw("bind request body param error", "error", err)
			c.JSON(http.StatusBadRequest, err)
			return nil
		}

		form.Guard(context.Background(), guardReq)

		return nil
	}
}

func ProxyPoly(poly permit.Poly) echo.HandlerFunc {
	return func(c echo.Context) error {
		guardReq := &permit.GuardReq{
			Request: c.Request(),
			Writer:  c.Response().Writer,
			Body:    map[string]interface{}{},
		}
		if err := bindParams(c, guardReq); err != nil {
			logger.Logger.Errorw("bind request body param error", "error", err)
			c.JSON(http.StatusBadRequest, err)
			return nil
		}

		poly.Defender(context.Background(), guardReq)

		return nil
	}
}

func bindParams(c echo.Context, i *permit.GuardReq) error {
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
