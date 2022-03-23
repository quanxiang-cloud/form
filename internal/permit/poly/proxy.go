package defender

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"git.internal.yunify.com/qxp/misc/logger"
	"github.com/quanxiang-cloud/form/internal/permit"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

type Proxy struct {
	url       *url.URL
	transport http.RoundTripper

	next permit.Permit
}

func NewProxy(conf *config.Config) (*Proxy, error) {
	url, err := url.Parse(conf.Endpoint.Form)
	if err != nil {
		return nil, err
	}

	return &Proxy{
		url: url,
		transport: &http.Transport{
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
		},
	}, nil
}

func (p *Proxy) Do(ctx context.Context, req *permit.Request) (*permit.Response, error) {
	proxy := httputil.NewSingleHostReverseProxy(p.url)
	proxy.Transport = p.transport

	r := req.Request
	r.Host = p.url.Host
	data, err := json.Marshal(req.Body)
	if err != nil {
		logger.Logger.Errorf("entity json marshal failed: %s", err.Error())
		return nil, err
	}

	r.Body = io.NopCloser(bytes.NewReader(data))
	r.ContentLength = int64(len(data))
	proxy.ServeHTTP(req.Writer, r)

	return &permit.Response{}, nil
}
