package lowcode

import (
	"github.com/quanxiang-cloud/form/internal/models"
	"net/http"

	"github.com/quanxiang-cloud/cabin/tailormade/client"
)

type Form interface {
	GetCacheMatchRole(userID, depID, appID string) (*getMatchRoleResp, error)
	GetRoleMatchPermit(roleID string) (*findPermitResp, error)
}

type form struct {
	client http.Client
}

func NewForm(conf client.Config) Form {
	return &form{
		client: client.New(conf),
	}
}

type getMatchRoleResp struct {
	RoleID string `json:"roleID"`
	Types  int    `json:"types"`
}

func (f *form) GetCacheMatchRole(userID, depID, appID string) (*getMatchRoleResp, error) {
	return nil, nil
}

type findPermitResp struct {
	List []*permit `json:"list"`
}
type permit struct {
	ID        string             `json:"id"`
	Path      string             `json:"path"`
	Params    models.FiledPermit `json:"params"`
	Response  models.FiledPermit `json:"response"`
	Condition *models.Condition  `json:"condition"`
}

func (f *form) GetRoleMatchPermit(roleID string) (*findPermitResp, error) {
	//
	return nil, nil
}
