package client

import (
	"context"
	"fmt"
	"github.com/quanxiang-cloud/cabin/tailormade/client"
	"net/http"
)

const (
	polyapiHost = "http://polyapi:9090/api/v1/polyapi/inner/regSwagger/system/app/form"
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
func (p *polyapi) RegSwagger(ctx context.Context, host, swag, appID, tableID, tableName string) (*RegSwaggerResp, error) {
	namespace := fmt.Sprintf("/system/app/%s/raw/inner/form/%s", appID, tableID)
	params := struct {
		NameSpace      string `json:"namespace"`
		Host           string `json:"host"`
		Version        string `json:"version"`
		Swagger        string `json:"swagger"`
		NamespaceTitle string `json:"autoCreateNamespaceTitle"`
	}{
		Host:           host,
		Version:        version,
		Swagger:        swag,
		NameSpace:      namespace,
		NamespaceTitle: tableName,
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
	RegSwagger(ctx context.Context, host, swag, appID, tableID, tableName string) (*RegSwaggerResp, error)
}
