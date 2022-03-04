package client

import (
	"context"
	"github.com/quanxiang-cloud/cabin/tailormade/client"
	"github.com/quanxiang-cloud/form/pkg/misc/config"

	"github.com/gin-gonic/gin"
	"net/http"
)

const (

	//appCenterHost = "http://127.0.0.1:9012/api/v1/app-center"
	appCenterHost = "http://app-center/api/v1/app-center"
	checkIsAdmin  = "/checkIsAdmin"
	addAppScope   = "/addAppScope"
)

// CheckAppAdmin CheckAppAdmin
type CheckAppAdmin struct {
	IsAdmin bool
}

type appCenter struct {
	client http.Client
}

func (a *appCenter) AddAppScope(ctx context.Context, appID string, scopes []string) (*AddResp, error) {
	params := struct {
		AppID  string   `json:"appID"`
		Scopes []string `json:"scopes"`
	}{
		AppID:  appID,
		Scopes: scopes,
	}
	resp := &AddResp{}
	err := client.POST(ctx, &a.client, appCenterHost+addAppScope, params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// NewAppCenter 生成对象
func NewAppCenter(conf client.Config) AppCenter {
	return &appCenter{
		client: client.New(conf),
	}
}

// AddResp AddResp
type AddResp struct {
}

// AppCenter 应用壳管理对外接口
type AppCenter interface {
	CheckIsAdmin(ctx context.Context, appID, userID string, isSuper bool) (CheckAppAdmin, error)
	AddAppScope(ctx context.Context, appID string, scopes []string) (*AddResp, error)
}

func (a *appCenter) CheckIsAdmin(ctx context.Context, appID, userID string, isSuper bool) (CheckAppAdmin, error) {
	params := struct {
		AppID   string `json:"appID"`
		UserID  string `json:"userID"`
		IsSuper bool   `json:"is_super"`
	}{
		AppID:   appID,
		UserID:  userID,
		IsSuper: isSuper,
	}

	IsAdmin := CheckAppAdmin{}
	err := client.POST(ctx, &a.client, appCenterHost+checkIsAdmin, params, &IsAdmin)
	return IsAdmin, err
}

// AppCenterClient 应用壳服务请求客户端
type AppCenterClient struct {
	AppCenter AppCenter
}

// NewAppCenterClient NewAppCenterClient
func NewAppCenterClient(c *config.Config) *AppCenterClient {
	return &AppCenterClient{
		//AppCenter: NewAppCenter(),
	}
}

// CheckIsAppAdmin CheckIsAppAdmin
func (a *AppCenterClient) CheckIsAppAdmin(c *gin.Context) {
	//ctx := logger.CTXTransfer(c)
	//profile := header2.GetProfile(c)
	//appID := c.Param("appID")
	//if appID == "" {
	//	c.AbortWithStatus(http.StatusNotFound)
	//	return
	//} else if appID == "dataset" || appID == "formula" {
	//	c.Next()
	//	return
	//}
	//isSuper := header2.GetRole(c).IsSuper()
	//isAdmin, err := a.AppCenter.CheckIsAdmin(ctx, appID, profile.UserID, isSuper)
	//if err != nil {
	//	c.AbortWithStatus(http.StatusInternalServerError)
	//	return
	//}
	//if !isAdmin.IsAdmin {
	//	c.AbortWithStatus(http.StatusForbidden)
	//	return
	//}
	//c.Next()
	return
}
