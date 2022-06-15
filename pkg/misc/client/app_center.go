package client

import (
	"context"
	"fmt"
	"github.com/quanxiang-cloud/cabin/tailormade/client"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
	"strings"

	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	appCenter    = "/api/v1/app-center"
	checkIsAdmin = "/checkIsAdmin"
	getOne       = "/getOne"
)

type appCenterAPI struct {
	conf   *config.Config
	client http.Client
}

// NewAppCenterAPI 生成对象
func NewAppCenterAPI(conf *config.Config) AppCenterAPI {
	return &appCenterAPI{
		conf:   conf,
		client: client.New(conf.InternalNet),
	}
}

// AppCenterAPI 应用壳管理对外接口
type AppCenterAPI interface {
	CheckIsAdmin(ctx context.Context, appID, userID string, isSuper bool) (*CheckAppAdminResp, error)
	GetOne(ctx context.Context, appID string) (*AppResp, error)
}

// CheckAppAdminResp CheckAppAdmin
type CheckAppAdminResp struct {
	IsAdmin bool
}

func (a *appCenterAPI) CheckIsAdmin(ctx context.Context, appID, userID string, isSuper bool) (*CheckAppAdminResp, error) {
	params := struct {
		AppID   string `json:"appID"`
		UserID  string `json:"userID"`
		IsSuper bool   `json:"is_super"`
	}{
		AppID:   appID,
		UserID:  userID,
		IsSuper: isSuper,
	}
	resp := &CheckAppAdminResp{}
	err := client.POST(ctx, &a.client, fmt.Sprintf("%s%s%s", a.conf.Endpoint.AppCenter, appCenter, checkIsAdmin), params, resp)
	return resp, err
}

// AppCenterClient 应用壳服务请求客户端
type appCenterClient struct {
	appCenterAPI AppCenterAPI
}

// NewAppCenterClient NewAppCenterClient
func NewAppCenterClient(c *config.Config) *appCenterClient {
	return &appCenterClient{
		appCenterAPI: NewAppCenterAPI(c),
	}
}

// NewAppCenterMockClient NewAppCenterMockClient
func NewAppCenterMockClient(c *config.Config) *appCenterClient {
	return &appCenterClient{
		appCenterAPI: NewAppCenterMock(),
	}
}

// CheckIsAppAdmin CheckIsAppAdmin
func (a *appCenterClient) CheckIsAppAdmin(c *gin.Context) {
	ctx := header.MutateContext(c)
	appID := c.Param("appID")
	if appID == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	resp, err := a.appCenterAPI.CheckIsAdmin(ctx, appID, c.GetHeader("User-Id"), isSuper(GetRole(c)))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !resp.IsAdmin {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Next()
	return
}

func GetRole(c *gin.Context) []string {
	roleStr := c.Request.Header.Get("Role")
	return strings.Split(roleStr, ",")
}

func isSuper(roles []string) bool {
	for _, role := range roles {
		if role == "super" {
			return true
		}
	}
	return false
}

type AppResp struct {
	Id      string `json:"id"`
	PerPoly bool   `json:"perPoly"`
}

func (a *appCenterAPI) GetOne(ctx context.Context, appID string) (*AppResp, error) {
	resp := &AppResp{}
	err := client.POST(ctx, &a.client, fmt.Sprintf("%s%s%s", a.conf.Endpoint.AppCenter, appCenter, getOne), struct {
		AppID string `json:"appID"`
	}{
		AppID: appID,
	}, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
