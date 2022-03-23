package defender

import (
	"context"

	"github.com/quanxiang-cloud/form/internal/permit"
	"github.com/quanxiang-cloud/form/internal/permit/treasure"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

type Auth struct {
	next permit.Permit
	auth *treasure.Auth
}

func NewAuth(conf *config.Config) (*Auth, error) {
	auth, err := treasure.NewAuth(conf)
	if err != nil {
		return nil, err
	}

	next, err := NewProxy(conf)
	if err != nil {
		return nil, err
	}

	return &Auth{
		auth: auth,
		next: next,
	}, nil
}

func (a *Auth) Do(ctx context.Context, req *permit.Request) (*permit.Response, error) {
	p, err := a.auth.Auth(ctx, req)
	if err != nil {
		return nil, err
	}

	if p == nil {
		return nil, nil
	}
	return a.next.Do(ctx, req)
}
