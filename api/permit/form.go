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
	"github.com/quanxiang-cloud/form/internal/auth/condition"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

const (
	_userID = "User-Id"
	_depID  = "Department-Id"
	_appID  = "appID"
	_action = "action"
)

type Form struct {
	url  *url.URL
	form auth.Auth
	cond *condition.Condition
}

func NewForm(conf *config.Config) (*Form, error) {
	url, err := url.ParseRequestURI(conf.Endpoint.Form)
	if err != nil {
		return nil, err
	}

	form, err := auth.NewFormAuth(conf)
	if err != nil {
		return nil, err
	}
	cond := condition.NewCondition()

	return &Form{
		url:  url,
		form: form,
		cond: cond,
	}, nil
}

func (f *Form) auth(next echo.HandlerFunc) echo.HandlerFunc {
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

		havePermit, err := f.form.Auth(c.Request().Context(), formAuthReq)
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

func (f *Form) condition(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// reqData, err := io.ReadAll(c.Request().Body)
		// if err != nil {
		// 	logger.Logger.Errorw("read request body error", "error", err)
		// 	return err
		// }

		// condReq := &condition.CondReq{
		// 	UserID:   c.Request().Header.Get(_userID),
		// 	BodyData: make(map[string]interface{}),
		// }

		// err = json.Unmarshal(reqData, &condReq.BodyData)
		// if err != nil {
		// 	return err
		// }

		// err = f.cond.Do(c.Request().Context(), condReq)
		// if err != nil {
		// 	return err
		// }

		// data, err := json.Marshal(condReq.BodyData)
		// if err != nil {
		// 	return err
		// }

		// c.Request().Body = io.NopCloser(bytes.NewReader(data))
		// c.Request().ContentLength = int64(len(data))
		return next(c)
	}
}

func (f *Form) proxy() echo.HandlerFunc {
	return func(c echo.Context) error {
		proxy := httputil.NewSingleHostReverseProxy(f.url)
		// proxy.Transport = transport
		proxy.ModifyResponse = func(resp *http.Response) error {
			return f.form.Filter(resp, c.Param(_action))
		}

		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			logger.Logger.Errorf("Got error while modifying response: %v \n", err)
			return
		}

		r := c.Request()
		r.Host = f.url.Host
		proxy.ServeHTTP(c.Response().Writer, r)
		return nil
	}
}

func (f *Form) Forward(c echo.Context) error {
	ReqParam := &auth.ReqParam{
		Path: c.Request().URL.Path,
	}

	if err := bindParams(c, ReqParam); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return nil
	}

	resp, err := f.form.Auth(c.Request().Context(), ReqParam)
	if err != nil {
		c.JSON(http.StatusForbidden, nil)
		return nil
	}

	data, err := json.Marshal(resp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return nil
	}
	c.Request().Body = io.NopCloser(bytes.NewReader(data))
	c.Request().ContentLength = int64(len(data))

	proxy := httputil.NewSingleHostReverseProxy(f.url)
	// proxy.Transport = transport
	proxy.ModifyResponse = func(resp *http.Response) error {
		return f.form.Filter(resp, c.Param(_action))
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		logger.Logger.Errorf("Got error while modifying response: %v \n", err)
	}

	r := c.Request()
	r.Host = f.url.Host
	proxy.ServeHTTP(c.Response().Writer, r)
	return nil
}

func bindParams(c echo.Context, i *auth.ReqParam) error {
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
