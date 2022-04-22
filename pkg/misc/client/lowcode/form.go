package lowcode

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/quanxiang-cloud/cabin/tailormade/client"
	"github.com/quanxiang-cloud/form/internal/models"
)

var (
	formHost        string
	getUserRoleURL  = "%s/api/v1/form/%s/internal/apiRole/userRole/get"
	getPermitURl    = "%s/api/v1/form/%s/internal/apiPermit/get"
	saveUserRoleURL = "%s/api/v1/form/%s/internal/apiRole/userRole/create"
)

func init() {
	formHost = os.Getenv("FORM_HOST")
	if formHost == "" {
		formHost = "http://localhost:8080"
	}
}

type Form struct {
	client http.Client
}

func NewForm(conf client.Config) *Form {
	return &Form{
		client: client.New(conf),
	}
}

type GetMatchRoleResp struct {
	RoleID string          `json:"id"`
	Types  models.RoleType `json:"type"`
}

func (f *Form) GetCacheMatchRole(ctx context.Context, userID, depID, appID string) (*GetMatchRoleResp, error) {
	resp := &GetMatchRoleResp{}
	getUserRoleURLs := fmt.Sprintf(getUserRoleURL, formHost, appID)
	err := client.POST(ctx, &f.client, getUserRoleURLs, struct {
		UserID string `json:"userID"`
		DepID  string `json:"depID"`
		AppID  string `json:"appID"`
	}{
		UserID: userID,
		DepID:  depID,
		AppID:  appID,
	}, resp)
	if err != nil {
		return nil, err
	}

	if resp.RoleID == "" {
		return nil, nil
	}

	return resp, nil
}

type FindPermitResp struct {
	ID          string             `json:"id"`
	RoleID      string             `json:"roleID"`
	Path        string             `json:"path"`
	Params      models.FiledPermit `json:"params"`
	Response    models.FiledPermit `json:"response"`
	Condition   models.Condition   `json:"condition"`
	Methods     string             `json:"methods"`
	ResponseAll bool               `json:"responseAll"`
	ParamsAll   bool               `json:"paramsAll"`
}

func (f *Form) GetPermit(ctx context.Context, appID, roleID, path, methods string) (*FindPermitResp, error) {
	resp := &FindPermitResp{}
	getPermitURls := fmt.Sprintf(getPermitURl, formHost, appID)
	err := client.POST(ctx, &f.client, getPermitURls, struct {
		RoleID string `json:"roleID"`
		Path   string `json:"path"`
		Method string `json:"method"`
	}{
		RoleID: roleID,
		Path:   path,
		Method: methods,
	}, resp)
	if err != nil {
		return nil, err
	}
	if resp.ID == "" {
		return nil, nil
	}

	return resp, nil
}

type SaveRoleUsersResp struct {
}

func (f *Form) saveRoleUsers(ctx context.Context, appID, roleID, userID string) (*SaveRoleUsersResp, error) {
	resp := &SaveRoleUsersResp{}
	saveUserRoleURLs := fmt.Sprintf(saveUserRoleURL, formHost, appID)
	err := client.POST(ctx, &f.client, saveUserRoleURLs, struct {
		RoleID string `json:"roleID"`
		UserID string `json:"userID"`
		AppID  string `json:"appID"`
	}{
		RoleID: roleID,
		UserID: userID,
		AppID:  appID,
	}, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
