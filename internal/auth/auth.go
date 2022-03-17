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
	Auth(context.Context, *ReqParam) (bool, error)
	Filter(*http.Response, string) error
}

type auth struct {
	redis   models.LimitsRepo
	lowcode lowcode.Form
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

type ReqParam struct {
	AppID  string `param:"appID" query:"appID" header:"appID" form:"appID" json:"appID" xml:"appID"`
	UserID string `param:"userID" query:"userID" header:"User-Id" form:"userID" json:"userID" xml:"userID"`
	DepID  string `param:"depID" query:"depID" header:"Department-Id" form:"depID" json:"depID" xml:"depID"`
	Path   string `param:"path" query:"path" header:"path" form:"path" json:"path" xml:"path"`
	Body   map[string]interface{}
}

type RespParam struct {
	Match  *models.PermitMatch
	Permit *consensus.Permit
}

func (a *auth) Auth(ctx context.Context, req *ReqParam) (*RespParam, error) {
	// get the role information owned by the user
	match, err := a.getCacheMatch(ctx, req)
	if err != nil || match == nil {
		return nil, err
	}

	if match.Types == models.InitType {
		return nil, nil
	}

	permits, err := a.getCachePermit(ctx, match.RoleID, req)
	if err != nil {
		return nil, err
	}

	return &RespParam{
		Match: match,
		Permit: &consensus.Permit{
			Params:    permits.Params,
			Response:  permits.Response,
			Condition: permits.Condition,
			Types:     match.Types,
		},
	}, nil
}

func (a *auth) getCacheMatch(ctx context.Context, req *ReqParam) (*models.PermitMatch, error) {
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

func (a *auth) getCachePermit(ctx context.Context, roleID string, req *ReqParam) (*models.Limits, error) {
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
