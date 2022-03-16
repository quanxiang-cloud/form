package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/quanxiang-cloud/cabin/logger"

	"github.com/quanxiang-cloud/form/internal/auth"
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

	poly, err := NewPoly(polyEndpoint)
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(middleware.Logger(), middleware.Recover())

	e.Any("*", poly.proxy(), poly.auth)

	e.Any("api/v1/form/:appID/table/:tableID/*", form.proxy(), form.auth)

	logger.Logger.Info("start...")
	// Start server
	e.Start(port)
}

const (
	_userID       = "User-Id"
	_userName     = "User-Name"
	_departmentID = "Department-Id"
	_appID        = "appID"
	_tableID      = "tableID"
	_action       = "action"
)

type Form struct {
	url *url.URL
	fa  auth.FormAuth
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

	return &Form{
		url: url,
		fa:  fa,
	}, nil
}

func (f *Form) auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		res, err := f.fa.Auth(c.Request().Context(), &auth.FormAuthReq{
			AppID:   c.Param(_appID),
			TableID: c.Param(_tableID),
			Path:    c.Request().URL.Path,
		})
		if err != nil {
			return err
		}

		if !res.IsPermit {
			return errors.New("no permit")
		}

		return next(c)
	}
}

func (f *Form) proxy() echo.HandlerFunc {
	return func(c echo.Context) error {
		proxy := httputil.NewSingleHostReverseProxy(f.url)
		proxy.Transport = transport
		proxy.ModifyResponse = func(resp *http.Response) error {
			return f.fa.Filter(resp)
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
	poly auth.PolyAuth
}

func NewPoly(endpoint string) (*Poly, error) {
	url, err := url.ParseRequestURI(endpoint)
	if err != nil {
		return nil, err
	}

	poly := auth.NewPolyAuth()
	return &Poly{
		url:  url,
		poly: poly,
	}, nil
}

func (p *Poly) auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("a ...interface{}")
		return next(c)
	}
}

func (p *Poly) proxy() echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}