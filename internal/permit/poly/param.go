package defender

import (
	"context"
	"strings"

	"github.com/quanxiang-cloud/form/internal/permit"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

type Param struct {
	next permit.Permit
}

func NewParam(config *config.Config) (*Param, error) {
	next, err := NewAuth(config)
	if err != nil {
		return nil, err
	}

	return &Param{
		next: next,
	}, nil
}

const (
	minAPILength = 8
)

func (p *Param) Do(ctx context.Context, req *permit.Request) (*permit.Response, error) {
	pathArr := strings.Split(req.Request.URL.Path, "/")

	if len(pathArr) > minAPILength {
		req.AppID = pathArr[minAPILength-1]
	}

	return p.next.Do(ctx, req)
}
