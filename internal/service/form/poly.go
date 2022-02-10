package form

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/quanxiang-cloud/cabin/tailormade/client"
	"github.com/quanxiang-cloud/form/internal/service/types"
)

var (
	polyHost              string
	polyProxyTimeout      time.Duration
	polyProxymaxIdleConns int
)

const (
	proxyPath    = "/api/v1/polyapi/request/%s"
	proxyAPIPath = "system/app/%s/poly/default/%s.p"
)

func init() {
	flag.StringVar(&polyHost, "poly-host", "http://ployapi", "poly api host. default http://polyapi")
	flag.DurationVar(&polyProxyTimeout, "poly-proxy-timeout", 20*time.Second, "poly porxy timeout.default 20s")
	flag.IntVar(&polyProxymaxIdleConns, "poly-proxy-max-idle", 30, "poly proxy max idle conns. default 30")
}

type Poly struct {
	client http.Client
}

func NewPoly() *Poly {
	return &Poly{
		client: client.New(client.Config{
			Timeout:      polyProxyTimeout,
			MaxIdleConns: polyProxymaxIdleConns,
		}),
	}
}

type ProxyReq struct {
	base
	Page  int64
	Size  int64
	Query types.Query

	Action string
}

type ProxyResp interface{}

func (p *Poly) Proxy(ctx context.Context, req *ProxyReq) (*ProxyResp, error) {
	// FIXME 页面访问权限 或者 按钮权限
	params := map[string]interface{}{
		"userID": req.UserID,
	}
	queryTOMap(params, req)

	resp := new(ProxyResp)
	err := client.POST(ctx, &p.client, fmt.Sprintf(polyHost+proxyPath, fmt.Sprintf(proxyAPIPath, req.AppID, req.Action)), params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func queryTOMap(dst map[string]interface{}, params *ProxyReq) {
	for name, value := range params.Query {
		dst[name] = value
	}

	if params.Page != 0 {
		dst["page"] = params.Page
	}
	if params.Size != 0 {
		dst["size"] = params.Size
	}
}
