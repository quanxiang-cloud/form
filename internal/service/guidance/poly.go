package guidance

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"time"

	"git.internal.yunify.com/qxp/misc/client"
	"github.com/quanxiang-cloud/form/internal/service/consensus"
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
	flag.StringVar(&polyHost, "poly-host", "http://192.168.200.20:9017", "poly api host. default http://polyapi")
	flag.DurationVar(&polyProxyTimeout, "poly-proxy-timeout", 20*time.Second, "poly porxy timeout.default 20s")
	flag.IntVar(&polyProxymaxIdleConns, "poly-proxy-max-idle", 30, "poly proxy max idle conns. default 30")
}

type poly struct {
	client http.Client
}

func newPoly() (Guidance, error) {
	return &poly{
		client: client.New(client.Config{
			Timeout:      polyProxyTimeout,
			MaxIdleConns: polyProxymaxIdleConns,
		}),
	}, nil
}

func (p *poly) Do(ctx context.Context, bus *consensus.Bus) (*consensus.Response, error) {
	return p.proxy(ctx, bus)
}

func (p *poly) proxy(ctx context.Context, bus *consensus.Bus) (*consensus.Response, error) {
	params := map[string]interface{}{
		"userID": bus.UserID,
	}
	busTOMap(params, bus)

	resp := new(consensus.Response)
	err := client.POST(ctx, &p.client, fmt.Sprintf(polyHost+proxyPath, fmt.Sprintf(proxyAPIPath, bus.AppID, bus.Method)), params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func busTOMap(dst map[string]interface{}, bus *consensus.Bus) {
	queryToMap(dst, bus.Query)

	if bus.Entity != nil {
		dst["entity"] = bus.Entity
	}

	if bus.Page != 0 {
		dst["page"] = bus.Page
	}
	if bus.Size != 0 {
		dst["size"] = bus.Size
	}
}

func queryToMap(dst map[string]interface{}, query types.Query) {
	qd := qd{}
	if val, ok := query["bool"].(map[string]interface{}); ok {
		qd.boolSet(dst, val)
	}
}

// bool must must_not should match term terms range lt gt lte gte
type qd struct{}

func (qd *qd) boolSet(dst, src map[string]interface{}) {
	for key, val := range src {
		val, ok := val.([]interface{})
		if !ok {
			continue
		}
		switch key {
		case "must", "must_not", "should":
			for _, v := range val {
				et, ok := v.(map[string]interface{})
				if !ok {
					continue
				}
				qd.logicSet(dst, et)
			}
		default:
		}
	}
}

func (qd *qd) logicSet(dst, src map[string]interface{}) {
	for key, val := range src {
		val, ok := val.(map[string]interface{})
		if !ok {
			continue
		}

		if key == "bool" {
			qd.boolSet(dst, val)
			continue
		}

		qd.defaultSet(dst, val)
	}
}

func (qd *qd) defaultSet(dst, src map[string]interface{}) {
	for name, val := range src {
		dst[name] = val
	}
}
func (qd *qd) range1Set(dst, src map[string]interface{}) {
	// TODO
}
