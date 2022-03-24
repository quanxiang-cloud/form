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
	formHost           string
	cacheMatchRoleURL  = "%s/api/v1/form/%s/m/apiRole/userRoleMatch"
	roleMatchPermitURL = "%s/api/v1/form/%s/m/apiPermit/find"
)

func init() {
	formHost = os.Getenv("FORM_HOST")
	if formHost == "" {
		formHost = "http://127.0.0.1:81"
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
	RoleID string `json:"roleID"`
	Types  int    `json:"types"`
}

func (f *Form) GetCacheMatchRole(ctx context.Context, userID, depID, appID string) (*GetMatchRoleResp, error) {
	resp := &GetMatchRoleResp{}
	cacheMatchRoleURL := fmt.Sprintf(cacheMatchRoleURL, formHost, appID)
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
	Condition models.Condition   `json:"condition"`
}

func (f *Form) GetRoleMatchPermit(ctx context.Context, appID, roleID string) (*FindPermitResp, error) {
	resp := &FindPermitResp{}
	roleMatchPermitURL := fmt.Sprintf(roleMatchPermitURL, formHost, appID)
	err := client.POST(ctx, &f.client, roleMatchPermitURL, struct {
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
