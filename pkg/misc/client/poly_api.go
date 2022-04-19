package client

import (
	"context"
	"fmt"
	"github.com/quanxiang-cloud/cabin/tailormade/client"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
	"net/http"
)

const (
	register = "/api/v1/polyapi/inner/regSwagger/system/app/form"
	delete   = "/api/v1/polyapi/inner/deleteNamespace"
	version  = "last"
)

// NewPolyAPI 生成编排对象
func NewPolyAPI(conf *config.Config) PolyAPI {
	return &polyapi{
		client: client.New(conf.InternalNet),
		conf:   conf,
	}
}

type polyapi struct {
	client http.Client
	conf   *config.Config
}

type DeleteNamespaceResp struct {
}

func (p *polyapi) DeleteNamespace(ctx context.Context, appID, tableID string) (*DeleteNamespaceResp, error) {
	namespace := fmt.Sprintf("/system/app/%s/raw/inner/form/%s", appID, tableID)
	url := fmt.Sprintf("%s%s%s", p.conf.Endpoint.PolyInner, delete, namespace)
	params := struct{}{}
	resp := &DeleteNamespaceResp{}
	err := client.POST(ctx, &p.client, url, params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
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

	err := client.POST(ctx, &p.client, fmt.Sprintf("%s%s", p.conf.Endpoint.PolyInner, register), params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// PolyAPI PolyAPI
type PolyAPI interface {
	RegSwagger(ctx context.Context, host, swag, appID, tableID, tableName string) (*RegSwaggerResp, error)
	DeleteNamespace(ctx context.Context, appID, tableID string) (*DeleteNamespaceResp, error)
}
