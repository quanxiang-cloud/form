package httputil

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/form/internal/permit"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Transport Transport
func Transport(conf *config.Config) *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   conf.Transport.Timeout,
			KeepAlive: conf.Transport.KeepAlive,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          conf.Transport.MaxIdleConns,
		IdleConnTimeout:       conf.Transport.IdleConnTimeout,
		TLSHandshakeTimeout:   conf.Transport.TLSHandshakeTimeout,
		ExpectContinueTimeout: conf.Transport.ExpectContinueTimeout,
	}
}

type Proxys struct {
	Url       *url.URL
	Transport http.RoundTripper
}

func DoPoxy(ctx context.Context, req *permit.Request, p *Proxys, modify ModifyResponse) error {
	proxy := httputil.NewSingleHostReverseProxy(p.Url)
	proxy.Transport = p.Transport
	if modify != nil {
		proxy.ModifyResponse = modify
	}
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		logger.Logger.WithName("modify response").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	r := req.Echo.Request()
	r.Host = p.Url.Host
	if !IsQueryMethod(req.Echo.Request().Method) {
		data, err := json.Marshal(req.Data)
		if err != nil {
			logger.Logger.WithName("form proxy").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
			return err
		}
		r.Body = io.NopCloser(bytes.NewReader(data))
		r.ContentLength = int64(len(data))
	}
	proxy.ServeHTTP(req.Echo.Response(), r)
	return nil
}

type ModifyResponse func(*http.Response) error
