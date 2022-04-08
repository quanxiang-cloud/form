package guard

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/logger"
	cabinr "github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/form/internal/permit"
	"github.com/quanxiang-cloud/form/internal/service/consensus"
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
		logger.Logger.Errorf("Got error while modifying response: %v \n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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
	cabinResp := new(cabinr.Resp)
	conResp := &consensus.Response{}
	cabinResp.Data = conResp

	err = json.Unmarshal(respDate, cabinResp)
	if err != nil {
		return err
	}

	if cabinResp.Code != error2.Success {
		return nil
	}

	//switch req.Action {
	//case "get":
	//	treasure.Post(conResp.GetResp.Entity, req.Permit.Response)
	//case "search":
	//	treasure.Post(conResp.ListResp.Entities, req.Permit.Response)
	//}
	//
	data, err := json.Marshal(cabinResp)
	if err != nil {
		logger.Logger.Errorf("entity json marshal failed: %s", err.Error())
		return err
	}

	resp.Body = io.NopCloser(bytes.NewReader(data))
	resp.ContentLength = int64(len(data))
	return nil
}
