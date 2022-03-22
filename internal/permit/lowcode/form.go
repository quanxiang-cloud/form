package lowcode

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/quanxiang-cloud/form/internal/models"

	"github.com/quanxiang-cloud/cabin/tailormade/client"
)

var (
	cacheMatchRoleURL  = "%s/api/v1/form/permit/role/userRoleMatch"
	roleMatchPermitURL = "%s/api/v1/form/permit/apiPermit/find"
)

func init() {
	formHost := os.Getenv("FORM_HOST")
	if formHost == "" {
		formHost = "http://form"
	}
	cacheMatchRoleURL = fmt.Sprintf(cacheMatchRoleURL, formHost)
	roleMatchPermitURL = fmt.Sprintf(roleMatchPermitURL, formHost)
}

type Form interface {
	GetCacheMatchRole(context.Context, string, string, string) (*GetMatchRoleResp, error)
	GetRoleMatchPermit(context.Context, string) (*FindPermitResp, error)
}

type form struct {
	client http.Client
}

func NewForm(conf client.Config) Form {
	return &form{
		client: client.New(conf),
	}
}

type GetMatchRoleResp struct {
	RoleID string `json:"roleID"`
	Types  int    `json:"types"`
}

func (f *form) GetCacheMatchRole(ctx context.Context, userID, depID, appID string) (*GetMatchRoleResp, error) {
	var resp *GetMatchRoleResp
	err := client.POST(ctx, &f.client, cacheMatchRoleURL, struct {
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
	List []*permit `json:"list"`
}
type permit struct {
	ID        string             `json:"id"`
	Path      string             `json:"path"`
	Params    models.FiledPermit `json:"params"`
	Response  models.FiledPermit `json:"response"`
	Condition *models.Condition  `json:"condition"`
}

func (f *form) GetRoleMatchPermit(ctx context.Context, roleID string) (*FindPermitResp, error) {
	var resp *FindPermitResp
	err := client.POST(ctx, &f.client, cacheMatchRoleURL, struct {
		RoleID string `json:"roleID"`
	}{
		RoleID: roleID,
	}, resp)
	if err != nil {
		return nil, err
	}

	if len(resp.List) == 0 {
		return nil, nil
	}

	return resp, nil
}
