package client

import (
	"context"
	"fmt"
	"github.com/quanxiang-cloud/cabin/tailormade/client"
	"github.com/quanxiang-cloud/form/internal/service/tables/swagger"
	"net/http"
)

const (
	polyapiHost = "http://polyapi:9090/api/v1/polyapi/inner/regSwagger/system/app/" + swagger.Service
	version     = "last"
)

// NewPolyAPI 生成编排对象
func NewPolyAPI(conf client.Config) PolyAPI {
	return &polyapi{
		client: client.New(conf),
	}
}

type polyapi struct {
	client http.Client
}

// RegSwaggerResp RegSwaggerResp
type RegSwaggerResp struct {
}

// RegSwagger RegSwagger
func (p *polyapi) RegSwagger(ctx context.Context, host, swag, appID, contents string) (*RegSwaggerResp, error) {
	namespace := fmt.Sprintf("/system/app/%s/raw/inner/%s/%s", appID, swagger.NameSpace, contents)

	params := struct {
		NameSpace string `json:"namespace"`
		Host      string `json:"host"`
		Version   string `json:"version"`
		Swagger   string `json:"swagger"`
	}{

		Host:      host,
		Version:   version,
		Swagger:   swag,
		NameSpace: namespace,
	}
	resp := &RegSwaggerResp{}

	err := client.POST(ctx, &p.client, polyapiHost, params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// PolyAPI PolyAPI
type PolyAPI interface {
	RegSwagger(ctx context.Context, host, swag, appID, contents string) (*RegSwaggerResp, error)
}
