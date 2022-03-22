package permit

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/quanxiang-cloud/form/internal/service/consensus"
)

type Poly interface {
	Defender(context.Context, *GuardReq) (*GuardResp, error)
}

type Form interface {
	Guard(context.Context, *GuardReq) (*GuardResp, error)
}

type GuardReq struct {
	Request *http.Request
	Writer  http.ResponseWriter
	Header  Header
	Param   Param
	Get     Get
	Body    Query
	Permit  *consensus.Permit
}

type GuardResp struct{}

type Header struct {
	UserID string `header:"User-Id"`
	DepID  string `header:"Department-Id"`
}

type Param struct {
	AppID  string `param:"appID"`
	Action string `param:"action"`
}

type Get struct {
	Query     Query `query:"query"`
	Condition Query `query:"condition"`
	Entity    Query `query:"entity"`
}

type Query map[string]interface{}

func (q *Query) UnmarshalParam(param string) error {
	return json.Unmarshal([]byte(param), q)
}
