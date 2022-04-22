package guard

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/quanxiang-cloud/form/internal/permit/treasure"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
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
	proxy.ModifyResponse = func(resp *http.Response) error {
		return p.filter(resp, req)
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		logger.Logger.WithName("modify response").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	r := req.Echo.Request()
	r.Host = p.url.Host
	data, err := json.Marshal(req.Body)
	if err != nil {
		logger.Logger.WithName("form proxy").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		return nil, err
	}

	r.Body = io.NopCloser(bytes.NewReader(data))
	r.ContentLength = int64(len(data))
	proxy.ServeHTTP(req.Echo.Response(), r)

	return &permit.Response{}, nil
}

const (
	contentType         = "Content-Type"
	mimeApplicationJSON = "application/json"
)

func (p *Proxy) filter(resp *http.Response, req *permit.Request) error {
	if resp.StatusCode != http.StatusOK {
		return nil
	}

	ctype := resp.Header.Get(contentType)
	if !strings.HasPrefix(ctype, mimeApplicationJSON) {
		return fmt.Errorf("response data content-type is not %s", mimeApplicationJSON)
	}
	respDate, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var result map[string]interface{}

	if err := json.Unmarshal(respDate, &result); err != nil {
		return err
	}
	//if result["code"] != error2.Success {
	//	return nil
	//}
	switch req.Action {
	case "get", "search":
		if !req.Permit.ResponseAll {
			treasure.Filter(result, req.Permit.Response)
		}
	}
	data, err := json.Marshal(result)
	if err != nil {
		logger.Logger.Errorf("entity json marshal failed: %s", err.Error())
		return err
	}
	resp.Body = io.NopCloser(bytes.NewReader(data))
	//resp.ContentLength = int64(len(data))
	return nil
}
