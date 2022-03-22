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
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/internal/models/redis"
	"github.com/quanxiang-cloud/form/internal/permit"
	"github.com/quanxiang-cloud/form/internal/permit/condition"
	filters "github.com/quanxiang-cloud/form/internal/permit/filter"
	"github.com/quanxiang-cloud/form/internal/permit/lowcode"
	"github.com/quanxiang-cloud/form/internal/service/consensus"
	"github.com/quanxiang-cloud/form/pkg/misc/code"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

const (
	lockPermission = "lockPermission"
	lockPerMatch   = "lockPerMatch"
	lockTimeout    = time.Duration(30) * time.Second // 30秒
	timeSleep      = time.Millisecond * 500          // 0.5 秒
	_entity        = "entity"
)

type Auth struct {
	next    permit.Form
	redis   models.LimitsRepo
	lowcode lowcode.Form
}

func NewAuth(conf *config.Config) (*Auth, error) {
	redisClient, err := redis2.NewClient(conf.Redis)
	if err != nil {
		return nil, err
	}

	next, err := condition.NewCondition(conf)
	if err != nil {
		return nil, err
	}
	return &Auth{
		redis:   redis.NewLimitRepo(redisClient),
		lowcode: lowcode.NewForm(client.Config{}),
		next:    next,
	}, nil
}

func (a *Auth) Guard(ctx context.Context, req *permit.GuardReq) (*permit.GuardResp, error) {
	pass, err := a.auth(ctx, req)
	if err != nil {
		return nil, err
	}

	// no permit
	if !pass {
		return nil, nil
	}

	entity := req.Body[_entity]
	if req.Request.Method == http.MethodGet {
		entity = req.Get.Entity
	}

	// input parameter judgment
	if !filters.Pre(entity, req.Permit.Params) {
		return nil, nil
	}

	return a.next.Guard(ctx, req)
}

func (a *Auth) Defender(ctx context.Context, req *permit.GuardReq) (*permit.GuardResp, error) {
	pass, err := a.auth(ctx, req)
	if err != nil {
		return nil, err
	}

	if !pass {
		return nil, nil
	}

	return nil, nil
}

func (a *Auth) auth(ctx context.Context, req *permit.GuardReq) (bool, error) {
	// get the role information owned by the user
	match, err := a.getCacheMatch(ctx, req)
	if err != nil || match == nil {
		return false, err
	}

	if match.Types == models.InitType {
		return false, nil
	}

	permits, err := a.getCachePermit(ctx, match.RoleID, req)
	if err != nil {
		return false, err
	}

	req.Permit = &consensus.Permit{
		Params:    permits.Params,
		Response:  permits.Response,
		Condition: permits.Condition,
		Types:     match.Types,
	}

	return true, nil
}

func (a *Auth) getCacheMatch(ctx context.Context, req *permit.GuardReq) (*models.PermitMatch, error) {
	// relese lock
	defer a.redis.UnLock(ctx, lockPerMatch)
	for i := 0; i < 5; i++ {
		perMatch, err := a.redis.GetPerMatch(ctx, req.Header.UserID, req.Param.AppID)
		if err != nil {
			logger.Logger.Errorw(req.Header.UserID, header.GetRequestIDKV(ctx).Fuzzy(), err.Error())
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

	resp, err := a.lowcode.GetCacheMatchRole(ctx, req.Header.UserID, req.Header.DepID, req.Param.AppID)
	if err != nil || resp == nil {
		return nil, err
	}

	return &models.PermitMatch{
		RoleID: resp.RoleID,
		Types:  models.RoleType(resp.Types),
	}, nil
}

func (a *Auth) getCachePermit(ctx context.Context, roleID string, req *permit.GuardReq) (*models.Limits, error) {
	// relese lock
	defer a.redis.UnLock(ctx, lockPermission)
	for i := 0; i < 5; i++ {
		exist := a.redis.ExistsKey(ctx, roleID)
		if exist {
			// judge path
			getPermit, err := a.redis.GetPermit(ctx, roleID, req.Request.URL.Path)
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
		if value.Path == req.Request.URL.Path {
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
