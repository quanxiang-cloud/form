package treasure

import (
	"context"

	"time"

	"github.com/quanxiang-cloud/cabin/logger"

	error2 "github.com/quanxiang-cloud/cabin/error"
	redis2 "github.com/quanxiang-cloud/cabin/tailormade/db/redis"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/form/internal/models"
	"github.com/quanxiang-cloud/form/internal/models/redis"
	"github.com/quanxiang-cloud/form/internal/permit"
	"github.com/quanxiang-cloud/form/internal/service/consensus"
	"github.com/quanxiang-cloud/form/pkg/misc/client/lowcode"
	"github.com/quanxiang-cloud/form/pkg/misc/code"
	"github.com/quanxiang-cloud/form/pkg/misc/config"
)

const (
	lockPermission = "lockPermission"
	lockPerMatch   = "lockPerMatch"
	lockTimeout    = time.Duration(30) * time.Second // 30秒
	timeSleep      = time.Millisecond * 500          // 0.5 秒
)

type Auth struct {
	redis models.LimitsRepo
	form  *lowcode.Form
}

func NewAuth(conf *config.Config) (*Auth, error) {
	redisClient, err := redis2.NewClient(conf.Redis)
	if err != nil {
		return nil, err
	}

	return &Auth{
		redis: redis.NewLimitRepo(redisClient),
		form:  lowcode.NewForm(conf.InternalNet),
	}, nil
}

func (a *Auth) Auth(ctx context.Context, req *permit.Request) (*consensus.Permit, error) {
	match, err := a.getUserRole(ctx, req)
	if err != nil || match == nil {
		return nil, err
	}
	if match.RoleID == models.RoleInit {
		return nil, nil
	}
	permits, err := a.getCachePermit(ctx, match.RoleID, req)
	if err != nil || permits == nil {
		return nil, err
	}

	return &consensus.Permit{
		Params:    permits.Params,
		Response:  permits.Response,
		Condition: permits.Condition,
	}, nil
}

func (a *Auth) getUserRole(ctx context.Context, req *permit.Request) (*models.UserRoles, error) {
	resp, err := a.form.GetCacheMatchRole(ctx, req.UserID, req.DepID, req.AppID)
	if err != nil || resp == nil {
		return nil, err
	}
	if resp.Types == models.InitType {
		resp.RoleID = models.RoleInit
	}
	return &models.UserRoles{
		RoleID: resp.RoleID,
	}, nil
}

func (a *Auth) getCachePermit(ctx context.Context, roleID string, req *permit.Request) (*models.Limits, error) {
	resp, err := a.form.GetPermit(ctx, req.AppID, roleID)
	if err != nil || resp == nil {
		return nil, err
	}
	var getPermit *models.Limits
	for _, value := range resp.List {
		if value.Path == req.Request.URL.Path {
			per := &models.Limits{
				Path:      value.Path,
				Condition: value.Condition,
				Params:    value.Params,
				Response:  value.Response,
			}
			getPermit = per
		}
	}
	return getPermit, nil
}

func (a *Auth) getUserRole1(ctx context.Context, req *permit.Request) (*models.UserRoles, error) {
	for i := 0; i < 5; i++ {
		perMatch, err := a.redis.GetPerMatch(ctx, req.AppID, req.UserID)
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
	defer a.redis.UnLock(ctx, lockPerMatch)
	resp, err := a.form.GetCacheMatchRole(ctx, req.UserID, req.DepID, req.AppID)
	if err != nil || resp == nil {
		return nil, err
	}
	perMatch := &models.UserRoles{
		RoleID: resp.RoleID,
		UserID: req.UserID,
		AppID:  req.AppID,
	}
	if resp.Types == models.InitType {
		perMatch.RoleID = models.RoleInit
		resp.RoleID = models.RoleInit
	}
	err = a.redis.CreatePerMatch(ctx, perMatch)
	if err != nil {
		logger.Logger.Errorw("create per match")
	}

	// 入库

	return &models.UserRoles{
		RoleID: resp.RoleID,
	}, nil
}

func (a *Auth) getCachePermit1(ctx context.Context, roleID string, req *permit.Request) (*models.Limits, error) {
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
	// relese lock
	defer a.redis.UnLock(ctx, lockPermission)

	resp, err := a.form.GetPermit(ctx, req.AppID, roleID)
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
