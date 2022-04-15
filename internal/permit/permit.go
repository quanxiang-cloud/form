package permit

import (
	"context"
	"encoding/json"

	"github.com/labstack/echo/v4"
	"github.com/quanxiang-cloud/form/internal/service/consensus"
)

type Permit interface {
	Do(context.Context, *Request) (*Response, error)
}

type Request struct {
	Echo echo.Context
	// Request *http.Request
	// Writer  http.ResponseWriter

	Universal
	FormReq
}

type Response struct{}

type Universal struct {
	AppID  string `param:"appID"`
	UserID string `header:"User-Id"`
	DepID  string `header:"Department-Id"`
}

type FormReq struct {
	Action string `param:"action"`
	Body   Object
	Query  Object `query:"query"`
	Entity Object `query:"entity"`

	Permit *consensus.Permit
}

type Object map[string]interface{}

func (o *Object) UnmarshalParam(param string) error {
	return json.Unmarshal([]byte(param), o)
}
