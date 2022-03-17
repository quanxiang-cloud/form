package auth

import (
	"context"
	"net/http"

	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

type polyAuth struct {
	auth *auth
}

func NewPolyAuth(conf *config.Config) (Auth, error) {
	auth, err := newAuth(conf)
	return &polyAuth{
		auth: auth,
	}, err
}

func (p *polyAuth) Auth(ctx context.Context, req *ReqParam) (bool, error) {
	resp, err := p.auth.Auth(ctx, req)
	if err != nil {
		return false, err
	}

	if resp == nil {
		return false, nil
	}

	return true, nil
}

func (p *polyAuth) Filter(*http.Response, string) error {
	return nil
}
