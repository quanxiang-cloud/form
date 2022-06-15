package lowcode

import (
	"context"
	"fmt"
	"github.com/quanxiang-cloud/cabin/tailormade/client"
	"net/http"
	"os"
)

var (
	appHost   string
	getAppOne = "%s/api/v1/app-center/getOne"
)

func init() {
	appHost = os.Getenv("APP_CENTER_HOST")
	if appHost == "" {
		appHost = "http://app-center"
	}
}

type AppCenter struct {
	client http.Client
}

func NewAppCenter(conf client.Config) *AppCenter {
	return &AppCenter{
		client: client.New(conf),
	}
}

type AppResp struct {
	Id      string `json:"id"`
	PerPoly bool   `json:"perPoly"`
}

func (f *AppCenter) GetOne(ctx context.Context, appID string) (*AppResp, error) {
	resp := &AppResp{}
	urls := fmt.Sprintf(getAppOne, appHost)
	err := client.POST(ctx, &f.client, urls, struct {
		AppID string `json:"appID"`
	}{
		AppID: appID,
	}, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
