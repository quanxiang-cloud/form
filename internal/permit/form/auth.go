package guard

import (
	"context"

	"github.com/quanxiang-cloud/form/internal/permit"
	"github.com/quanxiang-cloud/form/internal/permit/treasure"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

type Auth struct {
	auth *treasure.Auth

	next permit.Permit
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

const (
	_entity = "entity"
)

func (a *Auth) Do(ctx context.Context, req *permit.Request) (*permit.Response, error) {
	p, err := a.auth.Auth(ctx, req)
	if err != nil || p == nil {
		return nil, err
	}
	//entity := req.Body[_entity]
	//if entity != nil {
	//	// input parameter judgment
	//	if !treasure.Pre(entity, p.Params) {
	//		return nil, nil
	//	}
	//}
	req.Permit = p
	return a.next.Do(ctx, req)
}
