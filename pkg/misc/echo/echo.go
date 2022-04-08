package echo

import (
	"context"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"go.uber.org/zap"
)

// predefined header.
const (
	RequestID = "Request-Id"
	Timezone  = "Timezone"
	TenantID  = "Tenant-Id"
)

// MutateContext mutate context.
func MutateContext(c echo.Context) context.Context {
	var (
		_requestID interface{} = "Request-Id"
		_timezone  interface{} = "Timezone"
		_tenantID  interface{} = "Tenant-Id"

		ctx = context.Background()
	)

	ctx = context.WithValue(ctx, _requestID, c.Request().Header.Get(RequestID))
	ctx = context.WithValue(ctx, _timezone, c.Request().Header.Get(Timezone))
	ctx = context.WithValue(ctx, _tenantID, c.Request().Header.Get(TenantID))

	return ctx
}

// GetRequestID get Request-Id from header.
func GetRequestID(c echo.Context) zap.Field {
	return ginHeader(c, header.RequestID)
}

// GetTimezone get Timezone from header.
func GetTimezone(c echo.Context) zap.Field {
	return ginHeader(c, header.Timezone)
}

func ginHeader(c echo.Context, key string) zap.Field {
	val := c.Request().Header.Get(key)
	return zap.String(key, val)
}

func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Start timer
		start := time.Now()
		path := c.Request().URL.Path
		raw := c.Request().URL.RawQuery

		// Process request
		err := next(c)

		if raw != "" {
			path = path + "?" + raw
		}

		// Stop timer
		logger.Logger.Infow("[Echo]",
			GetRequestID(c),
			zap.Float64("latency", time.Now().Sub(start).Seconds()),
			zap.String("clientIP", c.RealIP()),
			zap.String("method", c.Request().Method),
			zap.Int("statusCode", c.Response().Status),
			zap.String("errorMessage", err.Error()),
			zap.Int("bodySize", int(c.Response().Size)),
			zap.String("path", path),
		)

		return nil
	}
}

func Recover(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request(), false)
				if brokenPipe {
					logger.Logger.Error(c.Request().URL.Path,
						c.Request().Header.Get(RequestID),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					c.Error(err.(error))
					return
				}

				logger.Logger.Error("[Recovery from panic]",
					GetRequestID(c),
					zap.Any("error", err),
					zap.String("request", string(httpRequest)),
					zap.String("stack", string(debug.Stack())),
				)

				c.NoContent(http.StatusInternalServerError)
			}
		}()

		return next(c)
	}
}
