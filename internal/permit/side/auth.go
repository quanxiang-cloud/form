package side

import (
	"context"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/internal/permit"
	"github.com/quanxiang-cloud/form/internal/permit/treasure"
	httputil2 "github.com/quanxiang-cloud/form/pkg/httputil"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

// Auth is a guard for permit.
type Auth struct {
	auth *treasure.Auth
	next permit.Permit
}

// NewAuth returns a new guard for permit.
func NewAuth(conf *config.Config, rawurl string) (*Auth, error) {
	auth, err := treasure.NewAuth(conf)
	if err != nil {
		return nil, err
	}
	next, err := NewCondition(conf, rawurl)

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
	if p.Types == models.InitType {
		return a.next.Do(ctx, req)
	}
	if !p.ParamsAll {
		treasure.Filter(req.Data, p.Params)
	}
	if httputil2.IsQueryMethod(req.Echo.Request().Method) {
		req.Echo.Request().URL.RawQuery = httputil2.ObjectBodyToQuery(req.Data)
	}
	req.Permit = p
	return a.next.Do(ctx, req)
}
