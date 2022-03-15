package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/logger"
	authi "github.com/quanxiang-cloud/form/pkg/auth"
	"github.com/quanxiang-cloud/form/pkg/auth/lowcode"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

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

	conf, err := config.NewConfig(configPath)
	if err != nil {
		panic(err)
	}

	logger.Logger = logger.New(&conf.Log)
	if err != nil {
		panic(err)
	}

	formURI, err := url.ParseRequestURI(formEndpoint)
	if err != nil {
		panic(err)
	}

	polyURI, err := url.ParseRequestURI(polyEndpoint)
	if err != nil {
		panic(err)
	}

	permit := &permit{
		formURI: formURI,
		polyURI: polyURI,
		config:  conf,
		transport: &http.Transport{
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
		},
	}

	e := gin.New()
	e.Use(gin.Logger(), gin.Recovery())

	group := e.Group("/", permit.auth())
	group.Any("*path", permit.proxy())

	logger.Logger.Info("start...")
	e.Run(port)
}

const (
	_userID       = "User-Id"
	_userName     = "User-Name"
	_departmentID = "Department-Id"
	_appID        = "appID"
	_action       = "action"
)

type permit struct {
	formURI   *url.URL
	polyURI   *url.URL
	authi     authi.Interface
	transport *http.Transport
	config    *config.Config
}

func (p *permit) proxy() func(c *gin.Context) {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/v1/form") {
			p.serverHTTP(p.formURI, c)
		} else {
			p.serverHTTP(p.polyURI, c)
		}
	}
}

func (p *permit) auth() func(c *gin.Context) {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		data, err := c.GetRawData()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if strings.HasPrefix(path, "/api/v1/form") {

			req := &lowcode.FormReq{}
			req.UserID = c.GetHeader(_userID)
			req.UserName = c.GetHeader(_userName)
			req.Method = c.Param(_action)
			req.DepID = c.GetHeader(_departmentID)
			req.Path = c.Request.RequestURI
			req.AppID, req.TableID, err = checkURL(c)
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}

			err = json.Unmarshal(data, req)
			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}

			p.authi = lowcode.NewFormAuth(p.config, req)
		} else {
			p.authi = lowcode.NewPolyAuth()
		}

		if !p.authi.Auth(c) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(data))
		c.Next()
	}
}

func (p *permit) serverHTTP(url *url.URL, c *gin.Context) {
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.Transport = p.transport
	proxy.ModifyResponse = func(resp *http.Response) error {
		return p.authi.Filter(resp)
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		logger.Logger.Errorf("Got error while modifying response: %v \n", err)
		return
	}

	r := c.Request
	r.Host = url.Host
	proxy.ServeHTTP(c.Writer, r)
}

func checkURL(c *gin.Context) (appID, tableName string, err error) {
	appID, ok := c.Params.Get(_appID)
	tableName, okt := c.Params.Get("tableName")
	if !ok || !okt {
		err = errors.New("invalid URI")
		return
	}
	return
}
