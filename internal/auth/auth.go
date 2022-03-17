package auth

import (
	"context"
	"net/http"
	"time"

	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/cabin/tailormade/client"
	redis2 "github.com/quanxiang-cloud/cabin/tailormade/db/redis"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/form/internal/auth/filters"
	"github.com/quanxiang-cloud/form/internal/auth/lowcode"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/internal/models/redis"
	"github.com/quanxiang-cloud/form/internal/service/consensus"
	"github.com/quanxiang-cloud/form/pkg/misc/code"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

const (
	lockPermission = "lockPermission"
	lockPerMatch   = "lockPerMatch"
	lockTimeout    = time.Duration(30) * time.Second // 30秒
	timeSleep      = time.Millisecond * 500          // 0.5 秒
)

type Auth interface {
	Auth(context.Context, *AuthReq) (*AuthResp, error)
	Filter(*http.Response, string) error
}

type auth struct {
	redis   models.LimitsRepo
	lowcode lowcode.Form
	permit  *consensus.Permit
}

func newAuth(conf *config.Config) (*auth, error) {
	redisClient, err := redis2.NewClient(conf.Redis)
	if err != nil {
		return nil, err
	}

	return &auth{
		redis:   redis.NewLimitRepo(redisClient),
		lowcode: lowcode.NewForm(client.Config{}),
	}, nil
}

type AuthReq struct {
	AppID  string      `json:"appID,omitempty"`
	UserID string      `json:"userID,omitempty"`
	DepID  string      `json:"depID,omitempty"`
	Path   string      `json:"path,omitempty"`
	Entity interface{} `json:"entity,omitempty"`
}

type AuthResp struct {
	IsPermit bool `json:"isPermit,omitempty"`
}

func (a *auth) Auth(ctx context.Context, req *AuthReq) (*AuthResp, error) {
	return &AuthResp{true}, nil
	// get the role information owned by the user
	match, err := a.getCacheMatch(ctx, req)
	if err != nil || match == nil {
		return &AuthResp{}, err
	}

	if match.Types == models.InitType {
		return &AuthResp{true}, nil
	}

	permits, err := a.getCachePermit(ctx, match.RoleID, req)
	if err != nil {
		return &AuthResp{}, err
	}

	// access judgment
	if !filters.Pre(req.Entity, permits.Params) {
		return &AuthResp{}, error2.New(code.ErrNotPermit)
	}

	a.permit = &consensus.Permit{
		Params:    permits.Params,
		Response:  permits.Response,
		Condition: permits.Condition,
		Types:     match.Types,
	}

	return &AuthResp{true}, nil
}

func (a *auth) getCacheMatch(ctx context.Context, req *AuthReq) (*models.PermitMatch, error) {
	// relese lock
	defer a.redis.UnLock(ctx, lockPerMatch)
	for i := 0; i < 5; i++ {
		perMatch, err := a.redis.GetPerMatch(ctx, req.UserID, req.AppID)
		if err != nil {
			logger.Logger.Errorw(req.UserID, header.GetRequestIDKV(ctx).Fuzzy(), err.Error())
			return nil, err
		}
		if perMatch != nil {
			return perMatch, nil
		}

		// acquire distributed locks
		lock, err := a.redis.Lock(ctx, lockPerMatch, 1, lockTimeout)
		if err != nil {
			return nil, err
		}
		if !lock {
			<-time.After(timeSleep)
			continue
		}
		break
	}

	resp, err := a.lowcode.GetCacheMatchRole(ctx, req.UserID, req.DepID, req.AppID)
	if err != nil || resp == nil {
		return nil, err
	}

	return &models.PermitMatch{
		RoleID: resp.RoleID,
		Types:  models.RoleType(resp.Types),
	}, nil
}

func (a *auth) getCachePermit(ctx context.Context, roleID string, req *AuthReq) (*models.Limits, error) {
	// relese lock
	defer a.redis.UnLock(ctx, lockPermission)
	for i := 0; i < 5; i++ {
		exist := a.redis.ExistsKey(ctx, roleID)
		if exist {
			// judge path
			getPermit, err := a.redis.GetPermit(ctx, roleID, req.Path)
			if err != nil {
				return nil, err
			}
			if getPermit.Path == "" {
				return nil, error2.New(code.ErrNotPermit)
			}
			return getPermit, nil
		}

		// acquire distributed locks
		lock, err := a.redis.Lock(ctx, lockPermission, 1, lockTimeout)
		if err != nil {
			return nil, err
		}
		if !lock {
			<-time.After(timeSleep)
			continue
		}
		break
	}

	resp, err := a.lowcode.GetRoleMatchPermit(ctx, roleID)
	if err != nil || resp == nil {
		return nil, err
	}

	limits := make([]*models.Limits, len(resp.List))
	var getPermit *models.Limits
	for index, value := range resp.List {
		per := &models.Limits{
			Path:      value.Path,
			Condition: value.Condition,
			Params:    value.Params,
			Response:  value.Response,
		}
		if value.Path == req.Path {
			getPermit = per
		}
		limits[index] = per
	}
	err = a.redis.CreatePermit(ctx, roleID, limits...)
	if err != nil {
		logger.Logger.Errorw("create permit err", roleID, err.Error())
	}

	if getPermit == nil {
		return nil, error2.New(code.ErrNotPermit)
	}

	return getPermit, nil
}
