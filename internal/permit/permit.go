package permit

import (
	"context"
	"net/http"

	"github.com/quanxiang-cloud/form/internal/service/consensus"
)

type Permit interface {
	Do(context.Context, *Request) (*Response, error)
}

type Request struct {
	Request *http.Request
	Writer  http.ResponseWriter

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
	Body   Body
	Permit *consensus.Permit
}

type Body map[string]interface{}
