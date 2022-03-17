package router

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/form/internal/auth"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

type Poly struct {
	url  *url.URL
	poly auth.Auth
}

func NewPoly(conf *config.Config) (*Poly, error) {
	url, err := url.ParseRequestURI(conf.Endpoint.Poly)
	if err != nil {
		return nil, err
	}

	poly, err := auth.NewPolyAuth(conf)
	return &Poly{
		url:  url,
		poly: poly,
	}, err
}

func (p *Poly) auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		reqData, err := io.ReadAll(c.Request().Body)
		if err != nil {
			logger.Logger.Errorw("read request body error", "error", err)
			return err
		}

		formAuthReq := &auth.ReqParam{
			AppID:  c.Param(_appID),
			UserID: c.Request().Header.Get(_userID),
			DepID:  c.Request().Header.Get(_depID),
			Path:   c.Request().URL.Path,
		}

		err = json.Unmarshal(reqData, formAuthReq)
		if err != nil {
			logger.Logger.Errorw("unmarshal request body error", "error", err)
			return err
		}

		havePermit, err := p.poly.Auth(c.Request().Context(), formAuthReq)
		if err != nil {
			logger.Logger.Errorw("auth error", "error", err)
			return err
		}

		if !havePermit {
			c.Response().Writer.WriteHeader(http.StatusForbidden)
			return nil
		}

		c.Request().Body = io.NopCloser(bytes.NewReader(reqData))
		c.Request().ContentLength = int64(len(reqData))

		return next(c)
	}
}

func (p *Poly) proxy() echo.HandlerFunc {
	return func(c echo.Context) error {
		proxy := httputil.NewSingleHostReverseProxy(p.url)
		// proxy.Transport = transport
		proxy.ModifyResponse = func(resp *http.Response) error {
			return p.poly.Filter(resp, c.Param(_action))
		}

		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			logger.Logger.Errorf("Got error while modifying response: %v \n", err)
			return
		}

		r := c.Request()
		r.Host = p.url.Host
		proxy.ServeHTTP(c.Response().Writer, r)
		return nil
	}
}
