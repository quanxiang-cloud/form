package guard

import (
	"context"
	"github.com/quanxiang-cloud/form/internal/models"

	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/form/internal/permit"
	"github.com/quanxiang-cloud/form/internal/permit/treasure"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

// Auth is a guard for permit.
type Auth struct {
	auth *treasure.Auth

	next permit.Permit
}

// NewAuth returns a new guard for permit.
func NewAuth(conf *config.Config) (*Auth, error) {
	auth, err := treasure.NewAuth(conf)
	if err != nil {
		return nil, err
	}
	next, err := NewCondition(conf)
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
		logger.Logger.WithName("form auth").Errorw(err.Error(), header.GetRequestIDKV(ctx).Fuzzy()...)
		return nil, err
	}

	if p == nil {
		return nil, nil
	}
	req.Permit = p
	if p.Types == models.InitType {
		return a.next.Do(ctx, req)
	}
	switch req.Action {
	case "create", "update":
		if !p.ParamsAll {
			treasure.Filter(req.Body, p.Params)
		}
	}
	return a.next.Do(ctx, req)
}
