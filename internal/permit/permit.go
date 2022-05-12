package permit

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/quanxiang-cloud/form/internal/service/consensus"
)

type Permit interface {
	Do(context.Context, *Request) (*Response, error)
}

type Request struct {
	*http.Request
	*echo.Response

	Data map[string]interface{}
	Universal
	Permit *consensus.Permit
}

type Universal struct {
	AppID  string `param:"appID"`
	UserID string `header:"User-Id"`
	DepID  string `header:"Department-Id"`
}

type (
	Object   map[string]interface{}
	Response struct{}
)
