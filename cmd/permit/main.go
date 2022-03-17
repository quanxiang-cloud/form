package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/form/internal/auth"
	"github.com/quanxiang-cloud/form/internal/auth/condition"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

var transport *http.Transport

func main() {
	var port string
	var configPath string
	var formEndpoint string
	var polyEndpoint string
	var timeout time.Duration
	var keepAlive time.Duration
	var maxIdleConns int
	var idleConnTimeout time.Duration
	var tlsHandshakeTimeout time.Duration
	var expectContinueTimeout time.Duration

	flag.StringVar(&port, "port", ":40001", "service port default: :40001")
	flag.StringVar(&configPath, "config", "../../configs/config.yml", "profile address")
	flag.StringVar(&formEndpoint, "form-endpoint", "http://localhost:80", "service form endpoint default: http://localhost:80")
	flag.StringVar(&polyEndpoint, "poly-endpoint", "http://localhost:80", "service poly endpoint default: http://polyapi:80")
	flag.DurationVar(&timeout, "timeout", 20*time.Second, "Timeout is the maximum amount of time a dial will wait for a connect to complete. If Deadline is also set, it may fail earlier")
	flag.DurationVar(&keepAlive, "keep-alive", 20*time.Second, "KeepAlive specifies the interval between keep-alive probes for an active network connection.")
	flag.IntVar(&maxIdleConns, "max-idle-conns", 10, "MaxIdleConns controls the maximum number of idle (keep-alive) connections across all hosts. Zero means no limit.")
	flag.DurationVar(&idleConnTimeout, "idle-conn-timeout", 20*time.Second, "IdleConnTimeout is the maximum amount of time an idle (keep-alive) connection will remain idle before closing itself.")
	flag.DurationVar(&tlsHandshakeTimeout, "tls-handshake-timeout", 10*time.Second, "TLSHandshakeTimeout specifies the maximum amount of time waiting to wait for a TLS handshake. Zero means no timeout.")
	flag.DurationVar(&expectContinueTimeout, "expect-continue-timeout", 1*time.Second, "")
	flag.Parse()

	transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   timeout * time.Second,
			KeepAlive: keepAlive * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          maxIdleConns,
		IdleConnTimeout:       idleConnTimeout,
		TLSHandshakeTimeout:   tlsHandshakeTimeout,
		ExpectContinueTimeout: expectContinueTimeout * time.Second,
	}

	conf, err := config.NewConfig(configPath)
	if err != nil {
		panic(err)
	}

	logger.Logger = logger.New(&conf.Log)
	if err != nil {
		panic(err)
	}

	form, err := NewForm(formEndpoint, conf)
	if err != nil {
		panic(err)
	}

	poly, err := NewPoly(polyEndpoint, conf)
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(middleware.Logger(), middleware.Recover())

	e.Any("*", poly.proxy(), poly.auth)

	formG := e.Group("/api/v1/form")
	{
		formG.Any("/:appID/home/form/:tableID/:action", form.proxy(), form.condition)
	}

	e.Start(port)
}

const (
	_userID = "User-Id"
	_depID  = "Department-Id"
	_appID  = "appID"
	_action = "action"
)

type Form struct {
	url  *url.URL
	fa   auth.Auth
	cond *condition.Condition
}

func NewForm(endpoint string, conf *config.Config) (*Form, error) {
	url, err := url.ParseRequestURI(endpoint)
	if err != nil {
		return nil, err
	}

	fa, err := auth.NewFormAuth(conf)
	if err != nil {
		return nil, err
	}
	cond := condition.NewCondition()

	return &Form{
		url:  url,
		fa:   fa,
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

		havePermit, err := f.fa.Auth(c.Request().Context(), formAuthReq)
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
		reqData, err := io.ReadAll(c.Request().Body)
		if err != nil {
			logger.Logger.Errorw("read request body error", "error", err)
			return err
		}

		condReq := &condition.CondReq{
			UserID:   c.Request().Header.Get(_userID),
			BodyData: make(map[string]interface{}),
		}

		err = json.Unmarshal(reqData, &condReq.BodyData)
		if err != nil {
			return err
		}

		err = f.cond.Do(c.Request().Context(), condReq)
		if err != nil {
			return err
		}

		data, err := json.Marshal(condReq.BodyData)
		if err != nil {
			return err
		}

		c.Request().Body = io.NopCloser(bytes.NewReader(data))
		c.Request().ContentLength = int64(len(data))
		return next(c)
	}
}

func (f *Form) proxy() echo.HandlerFunc {
	return func(c echo.Context) error {
		proxy := httputil.NewSingleHostReverseProxy(f.url)
		proxy.Transport = transport
		proxy.ModifyResponse = func(resp *http.Response) error {
			return f.fa.Filter(resp, c.Param(_action))
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

type Poly struct {
	url  *url.URL
	poly auth.Auth
}

func NewPoly(endpoint string, conf *config.Config) (*Poly, error) {
	url, err := url.ParseRequestURI(endpoint)
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
		proxy.Transport = transport
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
