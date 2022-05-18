package permit

import (
	"context"
	"github.com/quanxiang-cloud/form/internal/service/consensus"

	"github.com/labstack/echo/v4"
)

type Permit interface {
	Do(context.Context, *Request) (*Response, error)
}

type Request struct {
	Echo echo.Context
	Data map[string]interface{}
	Universal
	Permit *consensus.Permit
	Path   string
}

type Universal struct {
	AppID  string `param:"appID"`
	UserID string `header:"User-Id"`
	DepID  string `header:"Department-Id"`
}

type Object map[string]interface{}
type Response struct{}
